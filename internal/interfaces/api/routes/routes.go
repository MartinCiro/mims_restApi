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
	"api_go/internal/interfaces/api/handlers/login"
	"api_go/internal/interfaces/api/middlewares"
)

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
// En la función setupAllRoutes, modificar:
func setupAllRoutes(mux *http.ServeMux, app *app.App) {
	// Inicializar middlewares
	authMiddleware := middlewares.NewAuthMiddleware(app.JWTService)

	// Inicializar handlers
	estadosHandler := estados.NewEstadosHandler(app.EstadoService)
	loginHandler := login.NewLoginHandler(app.LoginService)
	authHandler := auth.NewAuthHandler(app.AuthService) // NUEVO
	profileHandler := auth.NewProfileHandler(app.AuthService)

	// ========== RUTAS PÚBLICAS ==========

	// Health checks
	mux.HandleFunc("GET /{$}", common_handler.HealthHandler)
	mux.HandleFunc("GET /health", common_handler.HealthHandler)
	mux.HandleFunc("GET /ready", common_handler.ReadyHandler(app.DB, app.RedisCache))

	// Autenticación (públicas)
	mux.HandleFunc("POST /api/auth/login", loginHandler.Login)
	mux.Handle("POST /api/auth/register", authMiddleware.OptionalAuth(http.HandlerFunc(authHandler.Register))) // NUEVO

	// ========== RUTAS PROTEGIDAS ==========

	// Crear un subrouter para rutas protegidas
	protected := http.NewServeMux()

	// Rutas de Estados
	protected.HandleFunc("GET /api/estados", estadosHandler.ObtenerEstados)
	protected.HandleFunc("GET /api/estados/{id}", estadosHandler.ObtenerEstadoXid)
	protected.HandleFunc("POST /api/estados", estadosHandler.CrearEstado)
	protected.HandleFunc("PUT /api/estados/{id}", estadosHandler.ActualizarEstado)
	protected.HandleFunc("DELETE /api/estados/{id}", estadosHandler.EliminarEstado)

	// Ruta de perfil
	protected.HandleFunc("GET /api/profile", profileHandler.GetProfile)

	// Aplicar middleware de autenticación a las rutas protegidas
	mux.Handle("/api/estados", authMiddleware.Handler(protected))
	mux.Handle("/api/estados/", authMiddleware.Handler(protected))
	mux.Handle("/api/profile", authMiddleware.Handler(protected))
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
