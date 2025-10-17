package main

import (
	cfg "api_go/config"
	"api_go/core/common"
	"api_go/core/whatsapp"
	"api_go/infrastructure/database"
	"api_go/infrastructure/external"
	"api_go/infrastructure/repositories"
	api_whatsapp "api_go/interfaces/api/whatsapp"
	"log"

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
	whatsappRepo := repositories.NewGormWhatsAppRepository()
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

	// Endpoint de health check - ACTUALIZADO al formato ResponseBody
	router.GET("/health", func(c *gin.Context) {
		healthData := map[string]interface{}{
			"status":   "healthy",
			"service":  "whatsapp-api",
			"database": "connected",
		}

		response := common.NewResponseBody(true, 200, healthData)
		c.JSON(200, response)
	})

	// Endpoint de info - ACTUALIZADO
	router.GET("/info", func(c *gin.Context) {
		infoData := map[string]interface{}{
			"name":        "WhatsApp API",
			"version":     "1.0.0",
			"environment": cfg.Environment,
			"provider":    cfg.WhatsAppProvider,
		}

		response := common.NewResponseBody(true, 200, infoData)
		c.JSON(200, response)
	})

	// Iniciar servidor
	log.Println("🚀 Servidor iniciado en :8080")
	log.Println("💾 Base de datos: whatsapp_sessions.db")
	router.Run(":8080")
}
