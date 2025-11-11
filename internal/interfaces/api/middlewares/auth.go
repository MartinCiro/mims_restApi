package middlewares

import (
	"context"
	"net/http"
	"sync"

	"api_go/internal/core/auth"
	"api_go/internal/infrastructure/cookies"
	"api_go/internal/infrastructure/jwt"
	"api_go/internal/interfaces/api/common"
)

// Caché en memoria para almacenar información de usuarios autenticados
var (
	userCache = &sync.Map{} // Map thread-safe para concurrencia
)

type AuthMiddleware struct {
	jwtService   *jwt.JWTService
	authService  *auth.AuthService
	authRepo     auth.AuthPort
	cookieSigner *cookies.CookieSigner
}

func NewAuthMiddleware(jwtService *jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (am *AuthMiddleware) SetCookieSigner(cookieSigner *cookies.CookieSigner) {
	am.cookieSigner = cookieSigner
}

// SetAuthService establece el auth service (para OptionalAuth)
func (am *AuthMiddleware) SetAuthService(authService *auth.AuthService) {
	am.authService = authService
}

// Handler implementa el middleware de autenticación (requiere autenticación)
func (am *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var user *auth.User

		if am.cookieSigner != nil && am.authService != nil {
			cookie, err := r.Cookie("tk")
			if err == nil && cookie.Value != "" {
				user, err = am.authService.ValidateCookie(cookie.Value)
				if err == nil && user != nil {
					// Agregar usuario al contexto
					ctx = context.WithValue(ctx, "user", user)
					ctx = context.WithValue(ctx, "userID", user.ID)

					// Continuar con el siguiente handler
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				} else {
					cookies.ClearAuthCookie(w)
				}
			}
		}
		response := common.NewErrorResponse(401, "Debe iniciar sesión para continuar")
		common.WriteJSONResponse(w, response, 401)
	})
}

// OptionalAuth implementa el middleware de autenticación opcional
func (am *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if am.cookieSigner != nil && am.authService != nil {
			cookie, err := r.Cookie("tk")
			if err == nil && cookie.Value != "" {
				user, err := am.authService.ValidateCookie(cookie.Value)
				if err == nil && user != nil {
					ctx = context.WithValue(ctx, "user", user)
					ctx = context.WithValue(ctx, "userID", user.ID)
				} else {
					cookies.ClearAuthCookie(w)
				}
			}
		}

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

func (am *AuthMiddleware) RequireAuthAndPermission(permsMiddleware *PermissionsMiddleware, requiredPermissions []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return am.Handler(
			permsMiddleware.Handler(requiredPermissions)(next),
		)
	}
}
