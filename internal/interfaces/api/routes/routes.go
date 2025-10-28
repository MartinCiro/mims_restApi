// internal/interfaces/api/routes/routes.go
package routes

import (
	"log"
	"net/http"
	"os"
	"time"

	"api_go/config"
	"api_go/internal/app"
	"api_go/internal/interfaces/api/common"
	common_handler "api_go/internal/interfaces/api/handlers/common"
	"api_go/internal/interfaces/api/handlers/estados"
	"api_go/internal/interfaces/api/handlers/login"
	"api_go/internal/interfaces/api/middlewares"
)

// SetupRouter configura todas las rutas de la aplicación
func SetupRouter(app *app.App, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Servir archivos estáticos (equivalente a ServeStaticModule)
	if cfg.IsDevelopment() {
		setupStaticFiles(mux)
	}

	// Rutas de API bajo /api/
	SetupAPIRoutes(mux, app)

	// Rutas públicas
	setupPublicRoutes(mux, app)

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
	mux.HandleFunc("GET /", common_handler.HealthHandler)
	mux.HandleFunc("GET /ready", common_handler.ReadyHandler(app.DB, app.RedisCache))

	// Autenticación (pública)
	// Login (público) - handler específico
	loginHandler := login.NewLoginHandler(app.LoginService)
	mux.HandleFunc("POST /api/auth/login", loginHandler.Login)

	// Otras rutas públicas de auth
	//authHandler := auth.NewAuthHandler(app.AuthService)
	/* mux.HandleFunc("POST /api/auth/refresh", authHandler.RefreshToken)
	mux.HandleFunc("POST /api/auth/validate", authHandler.ValidateToken) */
}

// setupAPIRoutes configura rutas protegidas de la API
func SetupAPIRoutes(mux *http.ServeMux, app *app.App) {
	// Inicializar middlewares
	authMiddleware := middlewares.NewAuthMiddleware(app.JWTService)

	// Inicializar handlers
	estadosHandler := estados.NewEstadosHandler(app.EstadoService)

	// Grupo de rutas protegidas
	protected := http.NewServeMux()

	// Rutas de Estados
	protected.HandleFunc("GET /estados", estadosHandler.ObtenerEstados)
	protected.HandleFunc("GET /estados/{id}", estadosHandler.ObtenerEstadoXid)
	protected.HandleFunc("POST /estados", estadosHandler.CrearEstado)
	protected.HandleFunc("PUT /estados/{id}", estadosHandler.ActualizarEstado)
	protected.HandleFunc("DELETE /estados/{id}", estadosHandler.EliminarEstado)

	// ✅ Aplicar middlewares a rutas protegidas - usar "/" como prefijo
	mux.Handle("/", authMiddleware.Handler(protected))
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
