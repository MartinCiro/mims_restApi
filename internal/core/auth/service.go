package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"api_go/config"
	"api_go/internal/infrastructure/jwt"
	"api_go/internal/infrastructure/redis"
	"api_go/internal/interfaces/api/common"
	"api_go/pkg/logger"
	"api_go/pkg/utils"
)

type AuthService struct {
	authPort        AuthPort
	redisService    *redis.Cache
	jwtService      *jwt.JWTService
	passwordService *utils.PasswordService
	userRepo        UserRepositoryPort
	rolRepo         RolRepositoryPort
	estadoRepo      EstadoRepositoryPort
	config          *config.Config
}

func NewAuthService(
	authPort AuthPort,
	redisService *redis.Cache,
	jwtService *jwt.JWTService,
	passwordService *utils.PasswordService,
	userRepo UserRepositoryPort,
	rolRepo RolRepositoryPort,
	estadoRepo EstadoRepositoryPort,
	config *config.Config,
) *AuthService {
	return &AuthService{
		authPort:        authPort,
		redisService:    redisService,
		jwtService:      jwtService,
		passwordService: passwordService,
		userRepo:        userRepo,
		rolRepo:         rolRepo,
		estadoRepo:      estadoRepo,
		config:          config,
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
	ID       int      `json:"id"`
	Nombre   string   `json:"nombre"`
	Rol      string   `json:"rol"`
	Permisos []string `json:"permisos,omitempty"`
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
	token, err := s.jwtService.GenerateJWT(jwt.JwtPayload{
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
	// Convertir userID string a int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("ID de usuario inválido")
	}

	userCacheKey := fmt.Sprintf("user:%d", userIDInt)

	// Intentar obtener del cache primero
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
	var usuarioRetrieved *User
	if userData.IDUser == 0 {
		// Obtener usuario del puerto de autenticación (ahora incluye rol y permisos)
		usuarioRetrieved, err = s.authPort.RetrieveUserByID(ctx, userIDInt)
		if err != nil {
			return nil, fmt.Errorf("usuario no encontrado")
		}

		userData = UserCacheData{
			IDUser:   usuarioRetrieved.ID,
			Nombre:   usuarioRetrieved.Username,
			IDRol:    usuarioRetrieved.IDRol,
			Permisos: usuarioRetrieved.Permisos, // ✅ Ahora incluye permisos en cache
		}

		// Guardar en Redis para futuras consultas
		userDataJSON, err := json.Marshal(userData)
		if err != nil {
			log.Printf("Error marshaling user data for cache: %v", err)
		} else {
			if s.redisService != nil {
				err = s.redisService.Set(ctx, userCacheKey, string(userDataJSON), 24*time.Hour)
				if err != nil {
					log.Printf("Error setting user cache: %v", err)
				}
			}
		}
	} else {
		// Si usamos cache, obtener el usuario completo para el rol nombre
		usuarioRetrieved, err = s.authPort.RetrieveUserByID(ctx, userIDInt)
		if err != nil {
			return nil, fmt.Errorf("usuario no encontrado")
		}
	}

	// Construir respuesta del perfil con la información completa
	profileData := ProfileResponse{
		ID:       usuarioRetrieved.ID,
		Nombre:   usuarioRetrieved.Username,
		Rol:      usuarioRetrieved.RolNombre,
		Permisos: usuarioRetrieved.Permisos,
	}

	return &common.ResponseBody[ProfileResponse]{
		Success: true,
		Code:    200,
		Data:    profileData,
	}, nil
}

// RegisterRequest define la estructura para registro de usuarios
type RegisterRequest struct {
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	RolNombre *string `json:"rol_nombre,omitempty"` // Opcional para usuarios autenticados
}

// RegisterResponse define la respuesta del registro
type RegisterResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

// RegisterUser registra un nuevo usuario
func (s *AuthService) RegisterUser(ctx context.Context, req RegisterRequest, currentUser *User) (*RegisterResponse, error) {
	logger.Info("Iniciando registro de usuario", "username", req.Username, "email", req.Email)

	var rolNombre string
	var estadoNombre string = "activo"

	// Determinar rol según el contexto
	if currentUser != nil {
		// Usuario autenticado: validar permisos
		hasPermission := s.hasPermission(currentUser.Permisos, "usuarios.crear")
		if !hasPermission {
			return nil, fmt.Errorf("no tiene permisos para crear usuarios")
		}

		// Usar rol proporcionado o "usuario" por defecto
		if req.RolNombre != nil {
			rolNombre = *req.RolNombre
		} else {
			rolNombre = "usuario"
		}
	} else {
		// Usuario guest: rol "invitado"
		rolNombre = "invitado"
	}

	// Buscar ID del rol por nombre
	rolID, err := s.rolRepo.FindRolIDByName(ctx, rolNombre)
	if err != nil {
		return nil, fmt.Errorf("error buscando rol '%s': %v", rolNombre, err)
	}
	if rolID == 0 {
		return nil, fmt.Errorf("rol '%s' no encontrado", rolNombre)
	}

	// Buscar ID del estado "activo"
	estadoID, err := s.estadoRepo.FindEstadoIDByName(ctx, estadoNombre)
	if err != nil {
		return nil, fmt.Errorf("error buscando estado '%s': %v", estadoNombre, err)
	}
	if estadoID == 0 {
		return nil, fmt.Errorf("estado '%s' no encontrado", estadoNombre)
	}

	// Crear nuevo usuario
	newUser, err := NewUsuario(req.Username, req.Password, &rolID, &estadoID)
	if err != nil {
		return nil, err
	}

	// Guardar en base de datos
	userID, err := s.userRepo.CreateUser(ctx, newUser, req.Email)
	if err != nil {
		// Verificar si es error de duplicado
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, fmt.Errorf("El usuario ya existe, inicie sesión para continuar")
		}

		logger.Error("No se pudo crear el usuario", "error", err)

		// Para cualquier otro error
		return nil, fmt.Errorf("Ha ocurrido un error en el servidor, contacte al administrador")
	}

	logger.Info("✅ Usuario creado en BD", "userID", userID)

	// Obtener usuario completo con permisos
	user, err := s.authPort.RetrieveUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo usuario creado: %v", err)
	}

	// Generar token JWT
	token, err := s.jwtService.GenerateJWT(jwt.JwtPayload{
		IDUser:   user.ID,
		Username: user.Username,
		IDRol:    rolID,
	})
	if err != nil {
		return nil, fmt.Errorf("error generando token: %v", err)
	}

	logger.Info("✅ Registro completado exitosamente", "userID", userID, "username", user.Username)

	// Construir respuesta
	response := &RegisterResponse{
		Token:     token,
		ExpiresIn: s.config.JWTExpireTime,
	}

	return response, nil
}

// hasPermission verifica si el usuario tiene un permiso específico
func (s *AuthService) hasPermission(permisos []string, permission string) bool {
	for _, perm := range permisos {
		if perm == permission {
			return true
		}
	}
	return false
}
