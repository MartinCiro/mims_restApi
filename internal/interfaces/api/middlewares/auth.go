package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"api_go/internal/infrastructure/jwt"
	"api_go/internal/interfaces/api/common"
)

// Caché en memoria para almacenar información de usuarios autenticados
var (
	userCache = &sync.Map{} // Map thread-safe para concurrencia
)

type AuthMiddleware struct {
	jwtService *jwt.JWTService
}

func NewAuthMiddleware(jwtService *jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// Handler implementa el middleware de autenticación
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

		// Validar información del usuario
		if userInfo != nil && userInfo.UserInfo != nil && userInfo.UserInfo.UserInfo.Doc != "" {
			response := common.NewErrorResponse(401, "Token inválido")
			common.WriteJSONResponse(w, response, 401)
			return
		}

		// Almacenar usuario en caché (de forma segura con goroutines)
		if userInfo != nil && userInfo.UserInfo != nil {
			userID := userInfo.UserInfo.UserInfo.IDUser
			if userID != 0 {
				userCache.Store(userID, userInfo)
			}
		}

		// Adjuntar información del usuario al contexto
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", userInfo)

		/* // Si hay nuevo token, agregarlo al header de respuesta
		if newToken != "" {
			w.Header().Set("X-New-Token", newToken)
		} */

		// Continuar con el siguiente handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserInfo obtiene información del usuario desde la caché (método estático equivalente)
func GetUserInfo(idUser string) (interface{}, bool) {
	return userCache.Load(idUser)
}

// CleanUserCache limpia la caché de un usuario específico
func CleanUserCache(idUser string) {
	userCache.Delete(idUser)
}
