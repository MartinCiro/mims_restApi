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

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResult es un objeto del dominio, NO una respuesta HTTP
type LoginResult struct {
	Token   string `json:"token"`
	Usuario struct {
		ID       int      `json:"id"`
		Nombre   string   `json:"nombre"`
		Rol      *int     `json:"rol"`
		Permisos []string `json:"permisos,omitempty"`
	} `json:"usuario"`
}

// Execute retorna objetos del dominio, NO estructuras HTTP
func (s *LoginService) Execute(ctx context.Context, credentials LoginCredentials) (*LoginResult, error) {
	fmt.Printf("🔍 Iniciando proceso de login para: %s\n", credentials.Username)

	// 1. Validar credenciales con AuthPort
	user, err := s.authPort.RetrieveUser(ctx, auth.AuthData{
		Username: credentials.Username,
	})

	if err != nil {
		fmt.Printf("❌ Error en RetrieveUser: %v\n", err)
		return nil, fmt.Errorf("credenciales inválidas")
	}

	if user == nil {
		fmt.Printf("❌ Usuario no encontrado: %s\n", credentials.Username)
		return nil, fmt.Errorf("credenciales inválidas")
	}

	fmt.Printf("✅ Usuario validado: ID=%d\n", user.ID)

	// 2. Verificar contraseña
	passwordMatch := s.passwordService.ComparePasswords(credentials.Password, user.PasswordHash)
	fmt.Printf("🔍 Comparación de contraseña: %t\n", passwordMatch)

	if !passwordMatch {
		fmt.Printf("❌ Contraseña incorrecta para: %s\n", credentials.Username)
		return nil, fmt.Errorf("credenciales inválidas")
	}

	fmt.Printf("✅ Credenciales válidas para: %s\n", credentials.Username)

	// 3. Generar token JWT
	token, err := s.jwtService.GenerateJWT(jwt.UserInfo{
		IDUser:   user.ID,
		Username: user.Username,
		IDRol:    *user.IDRol,
	})
	if err != nil {
		return nil, fmt.Errorf("error generando token: %v", err)
	}

	// 4. Guardar en cache
	userCacheKey := fmt.Sprintf("user:%d", user.ID)
	userData := map[string]interface{}{
		"id_user":  user.ID,
		"nombre":   user.Username,
		"id_rol":   user.IDRol,
		"permisos": user.Permisos,
	}

	userDataJSON, _ := json.Marshal(userData)
	s.redisService.Set(ctx, userCacheKey, string(userDataJSON), 24*time.Hour)

	// 5. Construir resultado del dominio
	result := &LoginResult{
		Token: token,
		Usuario: struct {
			ID       int      `json:"id"`
			Nombre   string   `json:"nombre"`
			Rol      *int     `json:"rol"`
			Permisos []string `json:"permisos,omitempty"`
		}{
			ID:       user.ID,
			Nombre:   user.Username,
			Rol:      user.IDRol,
			Permisos: user.Permisos,
		},
	}

	return result, nil
}
