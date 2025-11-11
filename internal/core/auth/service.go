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
	IDUser   string   `json:"id_user"`
	Nombre   string   `json:"nombre"`
	IDRol    *int     `json:"id_rol"`
	Permisos []string `json:"permisos"`
}

type ProfileResponse struct {
	ID       string   `json:"id"`
	Nombre   string   `json:"nombre"`
	Rol      string   `json:"rol"`
	Permisos []string `json:"permisos,omitempty"`
}

type AuthResponse struct {
	Message   string    `json:"message"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

func (s *AuthService) LoginUser(ctx context.Context, req LoginRequest) (*AuthResponse, string, time.Time, error) {
	if req.Email == "" {
		return nil, "", time.Time{}, fmt.Errorf("email o contraseña inválida")
	}

	usuarioRetrieved, err := s.authPort.RetrieveUser(ctx, AuthData{Email: req.Email})
	if err != nil {
		logger.Error("❌ Error buscando usuario", "email", req.Email, "error", err)
		return nil, "", time.Time{}, fmt.Errorf("email o contraseña inválida")
	}

	if usuarioRetrieved == nil || usuarioRetrieved.ID == "" {
		logger.Error("❌ Usuario no encontrado", "email", req.Email)
		return nil, "", time.Time{}, fmt.Errorf("email o contraseña inválida")
	}

	usuario := NewUsuarioFromEncrypted(
		usuarioRetrieved.Username,
		usuarioRetrieved.PasswordHash,
		usuarioRetrieved.IDRol,
		usuarioRetrieved.IDEstado,
	)

	isPasswordValid := usuario.ComparePassword(req.Password, s.passwordService)
	if !isPasswordValid {
		return nil, "", time.Time{}, fmt.Errorf("email o contraseña inválida")
	}

	// ✅ OBTENER PERMISOS PARA EL USUARIO
	permisos, err := s.authPort.GetPermissionsByRoleID(ctx, *usuarioRetrieved.IDRol)
	if err != nil {
		logger.Error("❌ Error obteniendo permisos durante login",
			"userID", usuarioRetrieved.ID,
			"roleID", *usuarioRetrieved.IDRol,
			"error", err)
		return nil, "", time.Time{}, fmt.Errorf("error obteniendo permisos del usuario: %v", err)
	}

	// ✅ VERIFICAR QUE EL USUARIO TENGA PERMISOS
	if len(permisos) == 0 {
		logger.Warn("⚠️ Usuario sin permisos asignados durante login",
			"userID", usuarioRetrieved.ID,
			"roleID", *usuarioRetrieved.IDRol)
		return nil, "", time.Time{}, fmt.Errorf("usuario sin permisos asignados")
	}

	userCacheKey := fmt.Sprintf("user:%s", usuarioRetrieved.ID)
	eventKey := fmt.Sprintf("event:usuario.logeado:%s", usuarioRetrieved.ID)

	cachedUser, err := s.redisService.Get(ctx, userCacheKey)
	var userData UserCacheData

	if err == nil && cachedUser != "" {
		err = json.Unmarshal([]byte(cachedUser), &userData)
		if err != nil {
			s.redisService.Delete(ctx, userCacheKey)
			userData = UserCacheData{}
		}
	}

	if userData.IDUser == "" {
		// ✅ GUARDAR PERMISOS EN EL CACHE
		userData = UserCacheData{
			IDUser:   usuarioRetrieved.ID,
			Nombre:   usuarioRetrieved.Username,
			IDRol:    usuarioRetrieved.IDRol,
			Permisos: permisos,
		}

		userDataJSON, err := json.Marshal(userData)
		if err != nil {
			logger.Error("Error marshaling user cache data:", err)
		} else {
			logger.Info("🔍 Guardando usuario con permisos en cache",
				"userID", usuarioRetrieved.ID,
				"permisosCount", len(permisos))

			err = s.redisService.Set(ctx, userCacheKey, userDataJSON, 24*time.Hour)
			if err != nil {
				logger.Error("Error setting user cache:", err)
			} else {
				logger.Info("✅ Usuario y permisos guardados en cache",
					"userID", usuarioRetrieved.ID,
					"permisos", permisos)
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

	eventExists, err := s.redisService.Get(ctx, eventKey)
	if err != nil || eventExists == "" {
		err = s.redisService.Set(ctx, eventKey, "true", 1*time.Hour)
		if err != nil {
			logger.Error("Error setting event cache:", err)
		}
	}

	response := &AuthResponse{
		Message:   "Login exitoso",
		ExpiresAt: cookieData.ExpiresAt,
	}

	logger.Info("✅ Login completado exitosamente",
		"userID", usuarioRetrieved.ID,
		"permisosCount", len(permisos))

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
			logger.Error("Error parsing cached user data:", err)
		}
	}

	var usuarioRetrieved *User
	if len(userData.IDUser) == 0 {
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

	response := &AuthResponse{
		Message:   "Login exitoso",
		ExpiresAt: cookieData.ExpiresAt,
	}

	return response, signedCookie, cookieData.ExpiresAt, nil
}

func (s *AuthService) ValidateCookie(cookieValue string) (*User, error) {
	cookieData, err := s.cookieSigner.Verify(cookieValue)
	if err != nil {
		return nil, fmt.Errorf("cookie inválida: %v", err)
	}

	ctx := context.Background()

	// ✅ OBTENER PERMISOS - ES OBLIGATORIO
	permisos, err := s.authPort.GetPermissionsByRoleID(ctx, cookieData.RoleID)
	if err != nil {
		logger.Error("❌ Error crítico obteniendo permisos para validación de cookie",
			"userID", cookieData.UserID,
			"roleID", cookieData.RoleID,
			"error", err)
		return nil, fmt.Errorf("error validando permisos del usuario: %v", err)
	}

	// ✅ VERIFICAR QUE HAYA PERMISOS
	if len(permisos) == 0 {
		logger.Warn("⚠️ Usuario sin permisos asignados",
			"userID", cookieData.UserID,
			"roleID", cookieData.RoleID)
		return nil, fmt.Errorf("usuario sin permisos asignados")
	}

	return &User{
		ID:        cookieData.UserID,
		Username:  cookieData.Username,
		RolNombre: cookieData.Role,
		IDRol:     &cookieData.RoleID,
		Permisos:  permisos, // ✅ PERMISOS OBLIGATORIOS
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
