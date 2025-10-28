// internal/core/login/service.go
package login

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"api_go/internal/core/auth"
	"api_go/internal/infrastructure/jwt"
	"api_go/internal/infrastructure/redis"
	"api_go/internal/interfaces/api/common"
	"api_go/pkg/utils"
)

type LoginService struct {
	authPort        auth.AuthPort
	jwtService      *jwt.JWTService
	redisService    *redis.Cache
	passwordService *utils.PasswordService
}

func NewLoginService(authPort auth.AuthPort, jwtService *jwt.JWTService, redisService *redis.Cache, passwordService *utils.PasswordService) *LoginService {
	return &LoginService{
		authPort:        authPort,
		jwtService:      jwtService,
		redisService:    redisService,
		passwordService: passwordService,
	}
}

// LoginCredentials DTO específico para login
type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse DTO específico para respuesta de login
type LoginResponse struct {
	Token   string `json:"token"`
	Usuario struct {
		ID       int    `json:"id"`
		Nombre   string `json:"nombre"`
		Rol      *int   `json:"rol"`
		Permisos []int  `json:"permisos,omitempty"`
	} `json:"usuario"`
}

func (s *LoginService) Execute(ctx context.Context, credentials LoginCredentials) (*common.ResponseBody[LoginResponse], error) {
	fmt.Printf("🔍 Executing login for user: %s\n", credentials.Username)
	// 1. Validar credenciales con AuthPort
	user, err := s.authPort.RetrieveUser(ctx, auth.AuthData{
		Username: credentials.Username,
	})
	if err != nil {
		return nil, fmt.Errorf("credenciales inválidas")
	}

	// 2. Verificar contraseña
	if !s.passwordService.ComparePasswords(credentials.Password, user.PasswordHash) {
		return nil, fmt.Errorf("credenciales inválidas")
	}

	// 3. Generar token JWT
	token, err := s.jwtService.GenerateJWT(jwt.UserInfo{
		IDUser:   user.ID,
		Username: user.Username,
		IDRol:    *user.IDRol,
	})
	if err != nil {
		return nil, fmt.Errorf("error generando token: %v", err)
	}

	// 4. Guardar en cache (si es necesario)
	userCacheKey := fmt.Sprintf("user:%d", user.ID)
	userData := map[string]interface{}{
		"id_user":  user.ID,
		"nombre":   user.Username,
		"id_rol":   user.IDRol,
		"permisos": user.Permisos,
	}

	userDataJSON, _ := json.Marshal(userData)
	s.redisService.Set(ctx, userCacheKey, string(userDataJSON), 24*time.Hour)

	// 5. Construir respuesta
	responseData := LoginResponse{
		Token: token,
		Usuario: struct {
			ID       int    `json:"id"`
			Nombre   string `json:"nombre"`
			Rol      *int   `json:"rol"`
			Permisos []int  `json:"permisos,omitempty"`
		}{
			ID:     user.ID,
			Nombre: user.Username,
			Rol:    user.IDRol,
		},
	}

	return &common.ResponseBody[LoginResponse]{
		Success: true,
		Code:    200,
		Data:    responseData,
	}, nil
}
