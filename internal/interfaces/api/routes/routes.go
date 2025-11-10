package routes

import (
	"log"
	"net/http"
	"os"
	"time"

	"api_go/config"
	"api_go/internal/app"
	"api_go/internal/interfaces/api/common"
	"api_go/internal/interfaces/api/handlers/auth"
	common_handler "api_go/internal/interfaces/api/handlers/common"
	"api_go/internal/interfaces/api/handlers/estados"
	"api_go/internal/interfaces/api/middlewares"
)

// responseWriter wrapper para interceptar status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// SetupRouter configura todas las rutas de la aplicación
func SetupRouter(app *app.App, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Servir archivos estáticos
	if cfg.IsDevelopment() {
		setupStaticFiles(mux)
	}

	// Configurar todas las rutas
	setupAllRoutes(mux, app)

	// Aplicar middlewares globales
	return withGlobalMiddleware(mux, cfg)
}

// setupStaticFiles configura archivos estáticos
func setupStaticFiles(mux *http.ServeMux) {
	publicDir := "./public"

	if _, err := os.Stat(publicDir); err == nil {
		fs := http.FileServer(http.Dir(publicDir))
		mux.Handle("/api-docs/", http.StripPrefix("/api-docs", fs))
		log.Println("📚 Serviendo documentación en /api-docs")
	}
}

// setupAllRoutes configura todas las rutas (públicas y protegidas)
func setupAllRoutes(mux *http.ServeMux, app *app.App) {
	// Inicializar middlewares
	authMiddleware := middlewares.NewAuthMiddleware(app.JWTService)

	// Configurar CookieSigner si está disponible
	if app.CookieSigner != nil {
		authMiddleware.SetCookieSigner(app.CookieSigner)
	}

	// Configurar auth service en el middleware
	authMiddleware.SetAuthService(app.AuthService)

	// Inicializar permisos middleware
	permsMiddleware := middlewares.NewPermissionsMiddleware(app.RedisCache)

	// Inicializar handlers
	estadosHandler := estados.NewEstadosHandler(app.EstadoService)
	authHandler := auth.NewAuthHandler(app.AuthService)
	profileHandler := auth.NewProfileHandler(app.AuthService)

	// ========== RUTAS PÚBLICAS ==========

	// Health checks
	mux.HandleFunc("GET /{$}", common_handler.HealthHandler)
	mux.HandleFunc("GET /health", common_handler.HealthHandler)
	mux.HandleFunc("GET /ready", common_handler.ReadyHandler(app.DB, app.RedisCache))

	// Autenticación (públicas)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.Handle("POST /api/auth/register", authMiddleware.OptionalAuth(http.HandlerFunc(authHandler.Register)))

	// ========== RUTAS PROTEGIDAS ==========

	// Crear un subrouter para rutas protegidas
	protected := http.NewServeMux()

	// Rutas de Estados con permisos
	protected.Handle("GET /api/estados",
		authMiddleware.RequireAuthAndPermission(permsMiddleware, middlewares.PermissionEstadosListar)(
			http.HandlerFunc(estadosHandler.ObtenerEstados),
		))

	protected.Handle("POST /api/auth/logout",
		authMiddleware.RequireAuthAndPermission(permsMiddleware, middlewares.PermissionLoginLogout)(
			http.HandlerFunc(authHandler.Logout),
		))

	protected.Handle("GET /api/estados/{id}",
		authMiddleware.RequireAuthAndPermission(permsMiddleware, middlewares.PermissionEstadosVer)(
			http.HandlerFunc(estadosHandler.ObtenerEstadoXid),
		))

	protected.Handle("POST /api/estados",
		authMiddleware.RequireAuthAndPermission(permsMiddleware, middlewares.PermissionEstadosCrear)(
			http.HandlerFunc(estadosHandler.CrearEstado),
		))

	protected.Handle("PUT /api/estados/{id}",
		authMiddleware.RequireAuthAndPermission(permsMiddleware, middlewares.PermissionEstadosEditar)(
			http.HandlerFunc(estadosHandler.ActualizarEstado),
		))

	protected.Handle("DELETE /api/estados/{id}",
		authMiddleware.RequireAuthAndPermission(permsMiddleware, middlewares.PermissionEstadosEliminar)(
			http.HandlerFunc(estadosHandler.EliminarEstado),
		))

	protected.Handle("GET /api/profile",
		authMiddleware.Handler(http.HandlerFunc(profileHandler.GetProfile)))

	// Aplicar middleware de autenticación base a las rutas protegidas
	mux.Handle("/api/estados", authMiddleware.Handler(protected))
	mux.Handle("/api/estados/", authMiddleware.Handler(protected))
	mux.Handle("/api/auth/me", authMiddleware.Handler(protected))
	mux.Handle("/api/auth/logout", authMiddleware.Handler(protected))
	mux.Handle("/api/profile", authMiddleware.Handler(protected))
}

// withGlobalMiddleware aplica middlewares globales
func withGlobalMiddleware(handler http.Handler, cfg *config.Config) http.Handler {
	// CORS - importante para cookies
	handler = withCORS(handler)

	// Logging
	handler = withLogging(handler, cfg)

	// Recovery (manejo de panics)
	handler = withRecovery(handler)

	return handler
}

// withCORS middleware para CORS (actualizado para cookies)
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, jwt, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")    // IMPORTANTE para cookies
		w.Header().Set("Access-Control-Expose-Headers", "Set-Cookie") // Para que el frontend vea las cookies

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
