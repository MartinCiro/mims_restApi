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
	"api_go/internal/infrastructure/cookies"
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
	cookieSigner    *cookies.CookieSigner
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
	cookieSigner *cookies.CookieSigner,
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
		cookieSigner:    cookieSigner,
		passwordService: passwordService,
		userRepo:        userRepo,
		rolRepo:         rolRepo,
		estadoRepo:      estadoRepo,
		config:          config,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User *User `json:"user"`
}

type UserCacheData struct {
	IDUser   int      `json:"id_user"`
	Nombre   string   `json:"nombre"`
	IDRol    *int     `json:"id_rol"`
	Permisos []string `json:"permisos,omitempty"`
}

type ProfileResponse struct {
	ID       int      `json:"id"`
	Nombre   string   `json:"nombre"`
	Rol      string   `json:"rol"`
	Permisos []string `json:"permisos,omitempty"`
}

type AuthResponse struct {
	User      *User               `json:"user"`
	Token     *jwt.VerifyResponse `json:"token,omitempty"`
	ExpiresAt time.Time           `json:"expires_at,omitempty"`
}

func (s *AuthService) LoginUser(ctx context.Context, req LoginRequest) (*AuthResponse, string, time.Time, error) {
	if req.Email == "" {
		return nil, "", time.Time{}, fmt.Errorf("email o contraseña inválida")
	}

	logger.Info("🔍 Iniciando login", "email", req.Email)

	usuarioRetrieved, err := s.authPort.RetrieveUser(ctx, AuthData{Email: req.Email})
	if err != nil {
		logger.Error("❌ Error buscando usuario", "email", req.Email, "error", err)
		return nil, "", time.Time{}, fmt.Errorf("email o contraseña inválida")
	}

	if usuarioRetrieved == nil || usuarioRetrieved.ID == 0 {
		logger.Error("❌ Usuario no encontrado", "email", req.Email)
		return nil, "", time.Time{}, fmt.Errorf("email o contraseña inválida")
	}

	logger.Info("✅ Usuario encontrado", "userID", usuarioRetrieved.ID, "email", usuarioRetrieved.Email)

	usuario := NewUsuarioFromEncrypted(
		usuarioRetrieved.Username,
		usuarioRetrieved.PasswordHash,
		usuarioRetrieved.IDRol,
		usuarioRetrieved.IDEstado,
	)

	isPasswordValid := usuario.ComparePassword(req.Password, s.passwordService)
	if !isPasswordValid {
		logger.Error("❌ Contraseña inválida", "userID", usuarioRetrieved.ID)
		return nil, "", time.Time{}, fmt.Errorf("email o contraseña inválida")
	}

	logger.Info("✅ Contraseña válida", "userID", usuarioRetrieved.ID)

	userCacheKey := fmt.Sprintf("user:%d", usuarioRetrieved.ID)
	eventKey := fmt.Sprintf("event:usuario.logeado:%d", usuarioRetrieved.ID)

	cachedUser, err := s.redisService.Get(ctx, userCacheKey)
	var userData UserCacheData

	if err == nil && cachedUser != "" {
		err = json.Unmarshal([]byte(cachedUser), &userData)
		if err != nil {
			log.Printf("Error parsing cached user data: %v", err)
		}
	}

	if userData.IDUser == 0 {
		userData = UserCacheData{
			IDUser: usuarioRetrieved.ID,
			Nombre: usuarioRetrieved.Username,
			IDRol:  usuarioRetrieved.IDRol,
		}

		userDataWithPermissions := userData

		userDataJSON, err := json.Marshal(userDataWithPermissions)
		if err != nil {
			log.Printf("Error marshaling user data for cache: %v", err)
		} else {
			err = s.redisService.Set(ctx, userCacheKey, string(userDataJSON), 24*time.Hour)
			if err != nil {
				log.Printf("Error setting user cache: %v", err)
			}
		}
	}

	cookieData := &cookies.SignedCookieData{
		UserID:    usuarioRetrieved.ID,
		Username:  usuarioRetrieved.Username,
		Role:      usuarioRetrieved.RolNombre,
		RoleID:    *usuarioRetrieved.IDRol,
		ExpiresAt: time.Now().Add(time.Duration(s.config.JWTExpireTime) * time.Second),
		IssuedAt:  time.Now(),
	}

	signedCookie, err := s.cookieSigner.Sign(cookieData)
	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("error generando cookie de autenticación: %v", err)
	}

	jwtToken, _, err := s.jwtService.GenerateJWTWithExpiry(jwt.JwtPayload{
		IDUser:   usuarioRetrieved.ID,
		Username: usuarioRetrieved.Username,
		IDRol:    *usuarioRetrieved.IDRol,
	}, time.Duration(s.config.JWTExpireTime)*time.Second)

	var tokenResponse *jwt.VerifyResponse
	if err == nil {
		tokenResponse = &jwt.VerifyResponse{
			UserInfo: &jwt.JwtPayload{
				IDUser:   usuarioRetrieved.ID,
				Username: usuarioRetrieved.Username,
				IDRol:    *usuarioRetrieved.IDRol,
			},
			JWT: &jwtToken,
		}
	}

	eventExists, err := s.redisService.Get(ctx, eventKey)
	if err != nil || eventExists == "" {
		err = s.redisService.Set(ctx, eventKey, "true", 1*time.Hour)
		if err != nil {
			log.Printf("Error setting event cache: %v", err)
		}
	}

	response := &AuthResponse{
		User:      usuarioRetrieved,
		Token:     tokenResponse,
		ExpiresAt: cookieData.ExpiresAt,
	}

	return response, signedCookie, cookieData.ExpiresAt, nil
}

func (s *AuthService) GetUserProfile(ctx context.Context, userID string) (*common.ResponseBody[ProfileResponse], error) {
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("ID de usuario inválido")
	}

	userCacheKey := fmt.Sprintf("user:%d", userIDInt)

	cachedUser, err := s.redisService.Get(ctx, userCacheKey)
	var userData UserCacheData

	if err == nil && cachedUser != "" {
		err = json.Unmarshal([]byte(cachedUser), &userData)
		if err != nil {
			log.Printf("Error parsing cached user data: %v", err)
		}
	}

	var usuarioRetrieved *User
	if userData.IDUser == 0 {
		usuarioRetrieved, err = s.authPort.RetrieveUserByID(ctx, userIDInt)
		if err != nil {
			return nil, fmt.Errorf("usuario no encontrado")
		}

		userData = UserCacheData{
			IDUser:   usuarioRetrieved.ID,
			Nombre:   usuarioRetrieved.Username,
			IDRol:    usuarioRetrieved.IDRol,
			Permisos: usuarioRetrieved.Permisos,
		}

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
		usuarioRetrieved, err = s.authPort.RetrieveUserByID(ctx, userIDInt)
		if err != nil {
			return nil, fmt.Errorf("usuario no encontrado")
		}
	}

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

type RegisterRequest struct {
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	RolNombre *string `json:"rol_nombre,omitempty"`
}

func (s *AuthService) RegisterUser(ctx context.Context, req RegisterRequest, currentUser *User) (*AuthResponse, string, time.Time, error) {
	logger.Info("Iniciando registro de usuario", "username", req.Username, "email", req.Email)

	var rolNombre string
	var estadoNombre string = "activo"

	if currentUser != nil {
		hasPermission := s.hasPermission(currentUser.Permisos, "usuarios.crear")
		if !hasPermission {
			return nil, "", time.Time{}, fmt.Errorf("no tiene permisos para crear usuarios")
		}

		if req.RolNombre != nil {
			rolNombre = *req.RolNombre
		} else {
			rolNombre = "usuario"
		}
	} else {
		rolNombre = "invitado"
	}

	rolID, err := s.rolRepo.FindRolIDByName(ctx, rolNombre)
	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("error buscando rol '%s': %v", rolNombre, err)
	}
	if rolID == 0 {
		return nil, "", time.Time{}, fmt.Errorf("rol '%s' no encontrado", rolNombre)
	}

	estadoID, err := s.estadoRepo.FindEstadoIDByName(ctx, estadoNombre)
	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("error buscando estado '%s': %v", estadoNombre, err)
	}
	if estadoID == 0 {
		return nil, "", time.Time{}, fmt.Errorf("estado '%s' no encontrado", estadoNombre)
	}

	newUser, err := NewUsuario(req.Username, req.Password, &rolID, &estadoID)
	if err != nil {
		return nil, "", time.Time{}, err
	}

	userID, err := s.userRepo.CreateUser(ctx, newUser, req.Email)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, "", time.Time{}, fmt.Errorf("El usuario ya existe, inicie sesión para continuar")
		}

		logger.Error("No se pudo crear el usuario", "error", err)
		return nil, "", time.Time{}, fmt.Errorf("Ha ocurrido un error en el servidor, contacte al administrador")
	}

	logger.Info("✅ Usuario creado en BD", "userID", userID)

	user, err := s.authPort.RetrieveUserByID(ctx, userID)
	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("error obteniendo usuario creado: %v", err)
	}

	cookieData := &cookies.SignedCookieData{
		UserID:    user.ID,
		Username:  user.Username,
		Role:      rolNombre,
		RoleID:    rolID,
		ExpiresAt: time.Now().Add(time.Duration(s.config.JWTExpireTime) * time.Second),
		IssuedAt:  time.Now(),
	}

	signedCookie, err := s.cookieSigner.Sign(cookieData)
	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("error generando cookie de autenticación: %v", err)
	}

	logger.Info("✅ Registro completado exitosamente", "userID", userID, "username", user.Username)

	response := &AuthResponse{
		User:      user,
		ExpiresAt: cookieData.ExpiresAt,
	}

	return response, signedCookie, cookieData.ExpiresAt, nil
}

func (s *AuthService) ValidateCookie(cookieValue string) (*User, error) {
	cookieData, err := s.cookieSigner.Verify(cookieValue)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        cookieData.UserID,
		Username:  cookieData.Username,
		RolNombre: cookieData.Role,
		IDRol:     &cookieData.RoleID,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userID int) error {
	userCacheKey := fmt.Sprintf("user:%d", userID)
	if s.redisService != nil {
		err := s.redisService.Delete(ctx, userCacheKey)
		if err != nil {
			log.Printf("Error clearing user cache: %v", err)
		}
	}
	return nil
}

func (s *AuthService) hasPermission(permisos []string, permission string) bool {
	for _, perm := range permisos {
		if perm == permission {
			return true
		}
	}
	return false
}
