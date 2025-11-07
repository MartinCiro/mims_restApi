package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"api_go/internal/core/auth"
	"api_go/internal/infrastructure/jwt"
	"api_go/internal/interfaces/api/common"
)

// Caché en memoria para almacenar información de usuarios autenticados
var (
	userCache = &sync.Map{} // Map thread-safe para concurrencia
)

type AuthMiddleware struct {
	jwtService  *jwt.JWTService
	authService *auth.AuthService
	authRepo    auth.AuthPort
}

func NewAuthMiddleware(jwtService *jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// SetAuthService establece el auth service (para OptionalAuth)
func (am *AuthMiddleware) SetAuthService(authService *auth.AuthService) {
	am.authService = authService
}

// Handler implementa el middleware de autenticación (requiere autenticación)
func (am *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extraer token del header
		rawToken := r.Header.Get("Authorization")
		if rawToken == "" {
			rawToken = r.Header.Get("jwt")
		}

		if rawToken == "" {
			response := common.NewErrorResponse(401, "No se ha proporcionado token")
			common.WriteJSONResponse(w, response, 401)
			return
		}

		// Eliminar prefijo "Bearer " si está presente
		token := strings.TrimPrefix(rawToken, "Bearer ")

		// Validar formato del token con expresión regular
		jwtRegex := regexp.MustCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$`)
		if !jwtRegex.MatchString(token) {
			response := common.NewErrorResponse(401, "El token proporcionado no tiene un formato válido")
			common.WriteJSONResponse(w, response, 401)
			return
		}

		// Verificar JWT
		userInfo, err := am.jwtService.VerifyJWT(token)
		if err != nil {
			fmt.Printf("Error verifying JWT: %v\n", err)
			response := common.NewErrorResponse(401, "Token inválido o expirado")
			common.WriteJSONResponse(w, response, 401)
			return
		}

		// ✅ CORREGIDO: Validación correcta del userInfo
		if userInfo == nil || userInfo.UserInfo == nil || userInfo.UserInfo.IDUser == 0 {
			response := common.NewErrorResponse(401, "Token inválido: información de usuario faltante")
			common.WriteJSONResponse(w, response, 401)
			return
		}

		// Almacenar usuario en caché
		userID := userInfo.UserInfo.IDUser
		if userID != 0 {
			userCache.Store(userID, userInfo)
		}

		// ✅ CORREGIDO: Agregar userID al contexto (como string)
		ctx := r.Context()
		userIDStr := strconv.Itoa(userID)
		ctx = context.WithValue(ctx, "userID", userIDStr)

		// También mantener el userInfo completo por si se necesita
		ctx = context.WithValue(ctx, "userInfo", userInfo)

		fmt.Printf("✅ Contexto actualizado con userID: %s\n", userIDStr)

		// Continuar con el siguiente handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth implementa el middleware de autenticación opcional
func (am *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extraer token del header
		rawToken := r.Header.Get("Authorization")
		if rawToken == "" {
			rawToken = r.Header.Get("jwt")
		}

		if rawToken == "" {
			// No hay token, continuar sin usuario
			next.ServeHTTP(w, r)
			return
		}

		// Eliminar prefijo "Bearer " si está presente
		token := strings.TrimPrefix(rawToken, "Bearer ")
		if token == rawToken {
			// No es un token Bearer, continuar sin usuario
			next.ServeHTTP(w, r)
			return
		}

		// Validar formato del token con expresión regular
		jwtRegex := regexp.MustCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$`)
		if !jwtRegex.MatchString(token) {
			// Token con formato inválido, continuar sin usuario
			next.ServeHTTP(w, r)
			return
		}

		// Verificar JWT
		userInfo, err := am.jwtService.VerifyJWT(token)
		if err != nil {
			// Token inválido, continuar sin usuario
			fmt.Printf("⚠️ Token inválido en OptionalAuth: %v\n", err)
			next.ServeHTTP(w, r)
			return
		}

		// Validar información del usuario
		if userInfo == nil || userInfo.UserInfo == nil || userInfo.UserInfo.IDUser == 0 {
			// Información de usuario faltante, continuar sin usuario
			next.ServeHTTP(w, r)
			return
		}

		// Si tenemos authService, obtener usuario completo con permisos
		var currentUser *auth.User
		if am.authService != nil {
			user, err := am.authRepo.RetrieveUserByID(r.Context(), userInfo.UserInfo.IDUser)
			if err == nil && user != nil {
				currentUser = user
				fmt.Printf("✅ Usuario autenticado encontrado: ID=%d, Username=%s\n", user.ID, user.Username)
			}
		}

		// Almacenar usuario en caché
		userID := userInfo.UserInfo.IDUser
		if userID != 0 {
			userCache.Store(userID, userInfo)
		}

		// Agregar información al contexto
		ctx := r.Context()
		userIDStr := strconv.Itoa(userID)
		ctx = context.WithValue(ctx, "userID", userIDStr)
		ctx = context.WithValue(ctx, "userInfo", userInfo)

		// Agregar usuario completo al contexto si está disponible
		if currentUser != nil {
			ctx = context.WithValue(ctx, "user", currentUser)
			fmt.Printf("✅ Usuario agregado al contexto: %s\n", currentUser.Username)
		}

		fmt.Printf("✅ OptionalAuth - Contexto actualizado con userID: %s\n", userIDStr)

		// Continuar con el siguiente handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserInfo obtiene información del usuario desde la caché
func GetUserInfo(idUser string) (interface{}, bool) {
	return userCache.Load(idUser)
}

// CleanUserCache limpia la caché de un usuario específico
func CleanUserCache(idUser string) {
	userCache.Delete(idUser)
}
