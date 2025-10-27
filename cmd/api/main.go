// cmd/api/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api_go/internal/app"
	"api_go/internal/interfaces/api/handlers/common"
	"api_go/internal/interfaces/api/routes"
)

func main() {
	// Cargar configuración
	cfg := config.Load()

	// Inicializar aplicación
	application := app.NewApp()
	if err := application.Initialize(cfg); err != nil {
		log.Fatalf("❌ Error inicializando aplicación: %v", err)
	}
	defer application.Shutdown()

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      routes.SetupRouter(application, cfg),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor en goroutine
	go func() {
		log.Printf("🚀 Server is running on http://localhost:%s", cfg.Port)
		log.Printf("📝 Environment: %s", cfg.Env)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Error iniciando servidor: %v", err)
		}
	}()

	// Esperar señal de interrupción para graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("🛑 Recibida señal de apagado, cerrando servidor...")

	// Graceful shutdown con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("⚠️ Error durante shutdown: %v", err)
	} else {
		log.Println("✅ Servidor cerrado correctamente")
	}
}

// setupRouter configura todas las rutas
func setupRouter(app *app.App) http.Handler {
	mux := http.NewServeMux()

	// Rutas públicas
	mux.HandleFunc("/", common.HealthHandler)
	mux.HandleFunc("/health", common.HealthHandler)
	mux.HandleFunc("/ready", common.ReadyHandler(app.DB, app.RedisCache))

	// Configurar rutas de la API
	apiRoutes := routes.SetupAPIRoutes(app)
	mux.Handle("/api/", http.StripPrefix("/api", apiRoutes))

	// Middleware global (CORS, logging, etc.)
	return withGlobalMiddleware(mux)
}

// withGlobalMiddleware aplica middlewares globales
func withGlobalMiddleware(handler http.Handler) http.Handler {
	// CORS middleware
	handler = withCORS(handler)

	// Logging middleware
	handler = withLogging(handler)

	return handler
}

// withCORS habilita CORS (equivalente a app.enableCors())
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configurar headers CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, jwt")

		// Manejar preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// withLogging middleware para log de requests
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Interceptar response writer para obtener status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, rw.statusCode, duration)
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
