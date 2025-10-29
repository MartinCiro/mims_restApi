package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"api_go/internal/infrastructure/jwt"
	"api_go/internal/infrastructure/redis"
	"api_go/internal/interfaces/api/common"
	"api_go/pkg/utils"
)

type AuthService struct {
	authPort        AuthPort
	redisService    *redis.Cache
	jwtService      *jwt.JWTService
	passwordService *utils.PasswordService
}

func NewAuthService(
	authPort AuthPort,
	redisService *redis.Cache,
	jwtService *jwt.JWTService,
	passwordService *utils.PasswordService,
) *AuthService {
	return &AuthService{
		authPort:        authPort,
		redisService:    redisService,
		jwtService:      jwtService,
		passwordService: passwordService,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Usuario struct {
		ID     int    `json:"id"`
		Nombre string `json:"nombre"`
		Rol    *int   `json:"rol"`
	} `json:"usuario"`
}

// UserCacheData representa los datos del usuario en cache
type UserCacheData struct {
	IDUser   int      `json:"id_user"`
	Nombre   string   `json:"nombre"`
	IDRol    *int     `json:"id_rol"`
	Permisos []string `json:"permisos,omitempty"`
}

// ProfileResponse define la estructura de respuesta para el perfil
type ProfileResponse struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Rol    *int   `json:"rol"`
	// Agrega otros campos que necesites exponer
	// Email    string   `json:"email,omitempty"`
	// Permisos []string `json:"permisos,omitempty"`
}

func (s *AuthService) LoginUser(ctx context.Context, req LoginRequest) (*common.ResponseBody[LoginResponse], error) {
	// Obtener usuario del puerto de autenticación
	usuarioRetrieved, err := s.authPort.RetrieveUser(ctx, AuthData{Username: req.Username})
	if err != nil {
		return nil, fmt.Errorf("usuario o contraseña inválida")
	}

	// Verificar contraseña
	usuario := NewUsuarioFromEncrypted(
		usuarioRetrieved.Username,
		usuarioRetrieved.PasswordHash,
		usuarioRetrieved.IDRol,
		usuarioRetrieved.IDEstado,
	)

	isPasswordValid := usuario.ComparePassword(req.Password, s.passwordService)
	if !isPasswordValid {
		return nil, fmt.Errorf("usuario o contraseña inválida")
	}

	userCacheKey := fmt.Sprintf("user:%d", usuarioRetrieved.ID)
	eventKey := fmt.Sprintf("event:usuario.logeado:%d", usuarioRetrieved.ID)

	// Revisar si el usuario ya está en cache
	cachedUser, err := s.redisService.Get(ctx, userCacheKey)
	var userData UserCacheData

	if err == nil && cachedUser != "" {
		// Usuario en cache, parsear datos
		err = json.Unmarshal([]byte(cachedUser), &userData)
		if err != nil {
			log.Printf("Error parsing cached user data: %v", err)
		}
	}

	// Si no hay datos en cache o hubo error, crear nuevos
	if userData.IDUser == 0 {
		userData = UserCacheData{
			IDUser: usuarioRetrieved.ID,
			Nombre: usuarioRetrieved.Username,
			IDRol:  usuarioRetrieved.IDRol,
		}

		// Agregar permisos si están disponibles (adaptar según tu estructura)
		userDataWithPermissions := userData
		// userDataWithPermissions.Permisos = usuarioRetrieved.Permisos // Descomentar si tienes permisos

		// Guardar en Redis
		userDataJSON, err := json.Marshal(userDataWithPermissions)
		if err != nil {
			log.Printf("Error marshaling user data for cache: %v", err)
		} else {
			err = s.redisService.Set(ctx, userCacheKey, string(userDataJSON), 24*time.Hour) // Cache por 24 horas
			if err != nil {
				log.Printf("Error setting user cache: %v", err)
			}
		}
	}

	// Generar JWT
	token, err := s.jwtService.GenerateJWT(jwt.UserInfo{
		IDUser:   usuarioRetrieved.ID,
		Username: usuarioRetrieved.Username,
		IDRol:    *usuarioRetrieved.IDRol,
	})
	if err != nil {
		return nil, fmt.Errorf("error generando token: %v", err)
	}

	// Verificar si el evento ya se publicó en Redis
	eventExists, err := s.redisService.Get(ctx, eventKey)
	if err != nil || eventExists == "" {
		// Evento no existe, guardarlo
		err = s.redisService.Set(ctx, eventKey, "true", 1*time.Hour) // Evento por 1 hora
		if err != nil {
			log.Printf("Error setting event cache: %v", err)
		}
	}

	// Construir respuesta
	responseData := LoginResponse{
		Token: token,
		Usuario: struct {
			ID     int    `json:"id"`
			Nombre string `json:"nombre"`
			Rol    *int   `json:"rol"`
		}{
			ID:     userData.IDUser,
			Nombre: userData.Nombre,
			Rol:    userData.IDRol,
		},
	}

	return &common.ResponseBody[LoginResponse]{
		Success: true,
		Code:    200,
		Data:    responseData,
	}, nil
}

func (s *AuthService) GetUserProfile(ctx context.Context, userID string) (*common.ResponseBody[ProfileResponse], error) {
	// Convertir userID string a int (ya que en tu estructura usas int)
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		// Usar el mismo patrón que LoginUser - retornar error para que el handler lo maneje
		return nil, fmt.Errorf("ID de usuario inválido")
	}

	userCacheKey := fmt.Sprintf("user:%d", userIDInt)

	// Intentar obtener del cache primero (misma lógica que LoginUser)
	cachedUser, err := s.redisService.Get(ctx, userCacheKey)
	var userData UserCacheData

	if err == nil && cachedUser != "" {
		// Usuario en cache, parsear datos
		err = json.Unmarshal([]byte(cachedUser), &userData)
		if err != nil {
			log.Printf("Error parsing cached user data: %v", err)
		}
	}

	// Si no hay datos en cache, obtener de la base de datos
	if userData.IDUser == 0 {
		// Obtener usuario del puerto de autenticación
		usuarioRetrieved, err := s.authPort.RetrieveUserByID(ctx, userIDInt)
		if err != nil {
			return nil, fmt.Errorf("usuario no encontrado")
		}

		userData = UserCacheData{
			IDUser: usuarioRetrieved.ID,
			Nombre: usuarioRetrieved.Username,
			IDRol:  usuarioRetrieved.IDRol,
		}

		// Guardar en Redis para futuras consultas (misma lógica que LoginUser)
		userDataJSON, err := json.Marshal(userData)
		if err != nil {
			log.Printf("Error marshaling user data for cache: %v", err)
		} else {
			err = s.redisService.Set(ctx, userCacheKey, string(userDataJSON), 24*time.Hour)
			if err != nil {
				log.Printf("Error setting user cache: %v", err)
			}
		}
	}

	// Construir respuesta del perfil (mismo patrón que LoginResponse)
	profileData := ProfileResponse{
		ID:     userData.IDUser,
		Nombre: userData.Nombre,
		Rol:    userData.IDRol,
	}

	// Usar el mismo patrón de respuesta que LoginUser
	return &common.ResponseBody[ProfileResponse]{
		Success: true,
		Code:    200,
		Data:    profileData,
	}, nil
}
