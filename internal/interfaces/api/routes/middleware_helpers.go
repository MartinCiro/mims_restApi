package routes

import (
	"net/http"

	"api_go/internal/interfaces/api/middlewares"
)

// ProtectedWithPermissions crea una cadena de middlewares para rutas protegidas con permisos
func ProtectedWithPermissions(
	authMiddleware *middlewares.AuthMiddleware,
	permsMiddleware *middlewares.PermissionsMiddleware,
	handler http.HandlerFunc,
	requiredPermissions []string,
) http.HandlerFunc {
	if len(requiredPermissions) > 0 {
		// Con permisos requeridos: Auth → Permissions → Handler
		// Convertir http.Handler a http.HandlerFunc
		chainedHandler := authMiddleware.Handler(
			permsMiddleware.Handler(requiredPermissions)(handler),
		)
		return chainedHandler.ServeHTTP
	}

	// Solo autenticación: Auth → Handler
	chainedHandler := authMiddleware.Handler(handler)
	return chainedHandler.ServeHTTP
}

// Protected crea una cadena solo con autenticación (sin permisos)
func Protected(
	authMiddleware *middlewares.AuthMiddleware,
	handler http.HandlerFunc,
) http.HandlerFunc {
	chainedHandler := authMiddleware.Handler(handler)
	return chainedHandler.ServeHTTP
}
