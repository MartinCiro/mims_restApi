package main

import (
	"log"
	cfg "scrapper_go_email/config"
	"scrapper_go_email/core/whatsapp"
	"scrapper_go_email/infrastructure/database"
	"scrapper_go_email/infrastructure/external"
	"scrapper_go_email/infrastructure/repositories"
	api_whatsapp "scrapper_go_email/interfaces/api/whatsapp"

	"github.com/gin-gonic/gin"
)

func main() {
	// Conectar a la base de datos
	cfg, err := cfg.Load()
	if err != nil {
		log.Fatal("❌ Error cargando configuración:", err)
	}

	// Configurar modo de Gin
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if err := database.Connect(cfg); err != nil {
		log.Fatal("❌ Error conectando a la base de datos:", err)
	}

	// Ejecutar migraciones
	if err := database.RunMigrations(); err != nil {
		log.Fatal("❌ Error ejecutando migraciones:", err)
	}

	// Infraestructura
	whatsappRepo := repositories.NewGormWhatsAppRepository() // Cambiado a GORM
	whatsmeowAdapter := external.NewWhatsMeowAdapter()

	// Core
	whatsappService := whatsapp.NewWhatsAppService(whatsappRepo, whatsmeowAdapter)

	// Interfaces
	whatsappController := api_whatsapp.NewWhatsAppController(whatsappService)

	// Router
	router := gin.Default()
	apiGroup := router.Group("/api/v1")

	// Registrar rutas
	api_whatsapp.RegisterRoutes(apiGroup, whatsappController)

	// Endpoint de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "healthy",
			"service":  "whatsapp-api",
			"database": "connected",
		})
	})

	// Iniciar servidor
	log.Println("🚀 Servidor iniciado en :8080")
	log.Println("💾 Base de datos: whatsapp_sessions.db")
	router.Run(":8080")
}
