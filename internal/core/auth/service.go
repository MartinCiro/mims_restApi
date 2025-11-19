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
	"api_go/internal/infrastructure/redis"
	"api_go/internal/interfaces/api/common"
	"api_go/pkg/logger"
	"api_go/pkg/utils"
)

type AuthService struct {
	authPort        AuthPort
	redisService    *redis.Cache
	cookieSigner    *cookies.CookieSigner
	passwordService *utils.PasswordService
	userRepo        UserRepositoryPort
	rolRepo         RolRepositoryPort
	estadoRepo      EstadoRepositoryPort
	config          *config.Config
	refreshService  *RefreshService
}

func NewAuthService(
	authPort AuthPort,
	redisService *redis.Cache,
	cookieSigner *cookies.CookieSigner,
	passwordService *utils.PasswordService,
	userRepo UserRepositoryPort,
	rolRepo RolRepositoryPort,
	estadoRepo EstadoRepositoryPort,
	config *config.Config,
) *AuthService {

	// Configurar tiempos de sesión
	sessionTTL := time.Duration(config.SessionTTL) * time.Second
	refreshThreshold := time.Duration(config.RefreshThreshold) * time.Second

	refreshService := NewRefreshService(
		redisService,
		cookieSigner,
		sessionTTL,
		refreshThreshold,
	)

	return &AuthService{
		authPort:        authPort,
		redisService:    redisService,
		cookieSigner:    cookieSigner,
		passwordService: passwordService,
		userRepo:        userRepo,
		rolRepo:         rolRepo,
		estadoRepo:      estadoRepo,
		config:          config,
		refreshService:  refreshService, // ← INYECTADO
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

	// ✅ CREAR SESIÓN EN REDIS después del login exitoso
	userIDInt, err := strconv.Atoi(usuarioRetrieved.ID)
	if err != nil {
		logger.Error("❌ Error convirtiendo userID a int para sesión",
			"userID", usuarioRetrieved.ID,
			"error", err)
		// No retornar error, solo loggear para no afectar el login
	} else {
		err = s.CreateUserSession(ctx, userIDInt, usuarioRetrieved.Username, usuarioRetrieved.RolNombre, *usuarioRetrieved.IDRol)
		if err != nil {
			logger.Error("❌ Error creando sesión en Redis",
				"userID", userIDInt,
				"error", err)
			// No retornar error, solo loggear para no afectar el login
		} else {
			logger.Info("✅ Sesión creada exitosamente en Redis",
				"userID", userIDInt,
				"username", usuarioRetrieved.Username,
				"role", usuarioRetrieved.RolNombre)
		}
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
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RolId    *int   `json:"rol_nombre,omitempty"`
}

func (s *AuthService) RegisterUser(ctx context.Context, req RegisterRequest, currentUser *User) (interface{}, string, time.Time, error) {
	var rolNombre string
	var estadoNombre string = "activo"

	if currentUser != nil {
		hasPermission := s.hasPermission(currentUser.Permisos, "usuario:crear")
		if !hasPermission {
			return nil, "", time.Time{}, fmt.Errorf("No tiene permisos para crear usuarios")
		}

		if req.RolId != nil {
			// CORRECCIÓN: Obtener el objeto Rol y luego extraer el nombre
			rol, err := s.rolRepo.FindByID(ctx, *req.RolId)
			if err != nil {
				return nil, "", time.Time{}, fmt.Errorf("error buscando rol por ID: %v", err)
			}
			if rol == nil {
				return nil, "", time.Time{}, fmt.Errorf("rol con ID %d no encontrado", *req.RolId)
			}
			rolNombre = rol.NombreRol // Asumiento que el campo se llama NombreRol
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

	// CAMBIO AQUÍ: Retornar diferentes estructuras según el caso
	if currentUser != nil {
		// Retornar solo el string para "Usuario creado correctamente"
		return "Usuario creado correctamente", signedCookie, cookieData.ExpiresAt, nil
	} else {
		return "Usuario registrado con exito", signedCookie, cookieData.ExpiresAt, nil
	}
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
	// ✅ ELIMINAR SESIÓN DE REDIS PRIMERO
	err := s.DeleteUserSession(ctx, userID)
	if err != nil {
		logger.Error("❌ Error eliminando sesión de Redis",
			"userID", userID,
			"error", err)
		// Continuar con el logout aunque falle eliminar la sesión
	} else {
		logger.Info("✅ Sesión eliminada de Redis", "userID", userID)
	}

	// ✅ LIMPIAR CACHE DE USUARIO (código existente)
	userCacheKey := fmt.Sprintf("user:%d", userID)
	if s.redisService != nil {
		err = s.redisService.Delete(ctx, userCacheKey)
		if err != nil {
			logger.Error("Error clearing user cache:",
				"userID", userID,
				"error", err)
		} else {
			logger.Info("✅ Cache de usuario eliminado", "userID", userID)
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

// CheckAndRefreshSession verifica y refresca la sesión si es necesario
func (s *AuthService) CheckAndRefreshSession(ctx context.Context, userID int, username string, role string, roleID int, currentCookie string) (*RefreshResult, error) {
	return s.refreshService.CheckAndRefresh(ctx, userID, username, role, roleID, currentCookie)
}

// CheckAndRefreshSessionFromCookie verifica y refresca directamente desde la cookie
func (s *AuthService) CheckAndRefreshSessionFromCookie(ctx context.Context, currentCookie string) (*RefreshResult, error) {
	return s.refreshService.CheckAndRefreshFromCookie(ctx, currentCookie)
}

// CreateUserSession crea una sesión para el usuario después del login
func (s *AuthService) CreateUserSession(ctx context.Context, userID int, username string, role string, roleID int) error {
	return s.refreshService.CreateSession(ctx, userID, username, role, roleID)
}

// UpdateUserSessionAccess actualiza el último acceso del usuario
func (s *AuthService) UpdateUserSessionAccess(ctx context.Context, userID int) error {
	return s.refreshService.UpdateSessionAccess(ctx, userID)
}

// GetSessionInfo obtiene información de la sesión actual
func (s *AuthService) GetSessionInfo(ctx context.Context, userID int) (map[string]interface{}, error) {
	return s.refreshService.GetSessionInfo(ctx, userID)
}

// ForceRefreshSession fuerza el refresh de una sesión
func (s *AuthService) ForceRefreshSession(ctx context.Context, userID int, username string, role string, roleID int) (*RefreshResult, error) {
	return s.refreshService.ForceRefresh(ctx, userID, username, role, roleID)
}

// DeleteUserSession elimina la sesión del usuario (para logout)
func (s *AuthService) DeleteUserSession(ctx context.Context, userID int) error {
	return s.refreshService.DeleteSession(ctx, userID)
}

// UserSessionExists verifica si el usuario tiene sesión activa
func (s *AuthService) UserSessionExists(ctx context.Context, userID int) (bool, error) {
	return s.refreshService.SessionExists(ctx, userID)
}
