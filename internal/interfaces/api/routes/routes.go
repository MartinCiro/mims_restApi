// internal/interfaces/api/routes/routes.go
package routes

import (
	"log"
	"net/http"
	"os"
	"time"

	"api_go/config"
	"api_go/internal/app"
	"api_go/internal/interfaces/api/handlers/auth"
	"api_go/internal/interfaces/api/handlers/common"
	"api_go/internal/interfaces/api/handlers/estados"
	"api_go/internal/interfaces/api/handlers/permisos"
	"api_go/internal/interfaces/api/handlers/roles"
	"api_go/internal/interfaces/api/handlers/usuarios"
	"api_go/internal/interfaces/api/middlewares"
)

// SetupRouter configura todas las rutas de la aplicación
func SetupRouter(app *app.App, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Servir archivos estáticos (equivalente a ServeStaticModule)
	if cfg.IsDevelopment() {
		setupStaticFiles(mux)
	}

	// Rutas públicas
	setupPublicRoutes(mux, app)

	// Rutas de API
	setupAPIRoutes(mux, app)

	// Aplicar middlewares globales
	return withGlobalMiddleware(mux, cfg)
}

// setupStaticFiles configura archivos estáticos
func setupStaticFiles(mux *http.ServeMux) {
	publicDir := "./public"

	// Verificar si existe el directorio public
	if _, err := os.Stat(publicDir); err == nil {
		fs := http.FileServer(http.Dir(publicDir))
		mux.Handle("/api-docs/", http.StripPrefix("/api-docs", fs))
		log.Println("📚 Serviendo documentación en /api-docs")
	}
}

// setupPublicRoutes configura rutas públicas
func setupPublicRoutes(mux *http.ServeMux, app *app.App) {
	// Health checks
	mux.HandleFunc("GET /", common.HealthHandler)
	mux.HandleFunc("GET /health", common.HealthHandler)
	mux.HandleFunc("GET /ready", common.ReadyHandler(app.DB, app.RedisCache))

	// Autenticación (pública)
	authHandler := auth.NewAuthHandler(app.AuthService)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
}

// setupAPIRoutes configura rutas protegidas de la API
func setupAPIRoutes(mux *http.ServeMux, app *app.App) {
	// Inicializar middlewares
	authMiddleware := middlewares.NewAuthMiddleware(app.JWTService)
	permsMiddleware := middlewares.NewPermissionsMiddleware(app.RedisCache)

	// Inicializar handlers
	estadosHandler := estados.NewEstadosHandler(app.EstadoService)
	usuariosHandler := usuarios.NewUsuariosHandler(app.UsuarioService)
	rolesHandler := roles.NewRolesHandler(app.RolService)
	permisosHandler := permisos.NewPermisosHandler(app.PermisoService)

	// Grupo de rutas protegidas
	protected := http.NewServeMux()

	// Rutas de Estados
	protected.HandleFunc("GET /estados", estadosHandler.ObtenerEstados)
	protected.HandleFunc("GET /estados/{id}", estadosHandler.ObtenerEstadoXid)
	protected.HandleFunc("POST /estados", estadosHandler.CrearEstado)
	protected.HandleFunc("PUT /estados/{id}", estadosHandler.ActualizarEstado)
	protected.HandleFunc("DELETE /estados/{id}", estadosHandler.EliminarEstado)

	// Rutas de Usuarios
	protected.HandleFunc("GET /usuarios", usuariosHandler.ObtenerUsuarios)
	protected.HandleFunc("GET /usuarios/{id}", usuariosHandler.ObtenerUsuarioXid)
	protected.HandleFunc("POST /usuarios", usuariosHandler.CrearUsuario)
	protected.HandleFunc("PUT /usuarios/{id}", usuariosHandler.ActualizarUsuario)
	protected.HandleFunc("DELETE /usuarios/{id}", usuariosHandler.EliminarUsuario)

	// Rutas de Roles
	protected.HandleFunc("GET /roles", rolesHandler.ObtenerRoles)
	protected.HandleFunc("GET /roles/{id}", rolesHandler.ObtenerRolXid)
	protected.HandleFunc("POST /roles", rolesHandler.CrearRol)
	protected.HandleFunc("PUT /roles/{id}", rolesHandler.ActualizarRol)
	protected.HandleFunc("DELETE /roles/{id}", rolesHandler.EliminarRol)

	// Rutas de Permisos
	protected.HandleFunc("GET /permisos", permisosHandler.ObtenerPermisos)
	protected.HandleFunc("GET /permisos/{id}", permisosHandler.ObtenerPermisoXid)
	protected.HandleFunc("POST /permisos", permisosHandler.CrearPermiso)
	protected.HandleFunc("PUT /permisos/{id}", permisosHandler.ActualizarPermiso)
	protected.HandleFunc("DELETE /permisos/{id}", permisosHandler.EliminarPermiso)

	// Aplicar middlewares a rutas protegidas
	mux.Handle("/api/", authMiddleware.Handler(protected))
}

// withGlobalMiddleware aplica middlewares globales
func withGlobalMiddleware(handler http.Handler, cfg *config.Config) http.Handler {
	// CORS
	handler = withCORS(handler)

	// Logging
	handler = withLogging(handler, cfg)

	// Recovery (manejo de panics)
	handler = withRecovery(handler)

	return handler
}

// withCORS middleware para CORS
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, jwt, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// withLogging middleware para logging
func withLogging(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		if cfg.IsDevelopment() {
			log.Printf("%s %s %d %v", r.Method, r.URL.Path, rw.statusCode, duration)
		}
	})
}

// withRecovery middleware para recuperación de panics
func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("⚠️ Panic recovered: %v", err)
				response := common.NewErrorResponse(500, "Internal Server Error")
				common.WriteJSONResponse(w, response, 500)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// responseWriter wrapper para interceptar status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
