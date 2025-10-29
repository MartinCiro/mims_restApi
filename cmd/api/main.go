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
	"api_go/pkg/logger"
)

func main() {
	// Cargar configuración
	cfg := config.Load()

	if cfg.Env == "Production" {
		logger.SetupProduction()
	} else {
		logger.SetupDevelopment()
	}
	logger.Info("iniciando aplicación",
		"version", "1.0.0",
		"environment", cfg.Env,
		"server_address", http.LocalAddrContextKey)

	// Inicializar aplicación
	application := app.NewApp()
	if err := application.Initialize(cfg); err != nil {
		logger.Fatal("fallo al inicializar aplicación", "error", err)
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
		logger.Info("iniciando servidor HTTP", "address", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("❌ Error iniciando servidor: %v", err)
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
