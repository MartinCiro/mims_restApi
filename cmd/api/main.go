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

	"api_go/config"
	"api_go/internal/app"
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
		/* log.Printf("📝 Environment: %s", cfg.Env) */

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

// responseWriter wrapper para interceptar status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
