package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Servidor
	Port        int
	Debug       bool
	Environment string

	// Base de datos
	DatabaseURL     string
	DatabaseLogMode bool

	// WhatsApp
	WhatsAppProvider string
	SessionTimeout   int

	// Gmail (si lo necesitas)
	GmailEmail       string
	GmailAppPassword string

	// Webhooks
	WebhookBaseURL string
}

func Load() (*Config, error) {
	// Cargar archivo .env (opcional - ignora el error si no existe)
	godotenv.Load()

	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, err
	}

	debug, err := strconv.ParseBool(getEnv("DEBUG", "true"))
	if err != nil {
		return nil, err
	}

	sessionTimeout, err := strconv.Atoi(getEnv("SESSION_TIMEOUT", "300"))
	if err != nil {
		return nil, err
	}

	databaseLogMode, err := strconv.ParseBool(getEnv("DATABASE_LOG_MODE", "false"))
	if err != nil {
		return nil, err
	}

	return &Config{
		// Servidor
		Port:        port,
		Debug:       debug,
		Environment: getEnv("ENVIRONMENT", "development"),

		// Base de datos
		DatabaseURL:     getEnv("DATABASE_URL", "whatsapp_sessions.db"),
		DatabaseLogMode: databaseLogMode,

		// WhatsApp
		WhatsAppProvider: getEnv("WHATSAPP_PROVIDER", "whatsmeow"),
		SessionTimeout:   sessionTimeout,

		// Gmail
		GmailEmail:       getEnv("GMAIL_EMAIL", ""),
		GmailAppPassword: getEnv("GMAIL_APP_PASSWORD", ""),

		// Webhooks
		WebhookBaseURL: getEnv("WEBHOOK_BASE_URL", ""),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
