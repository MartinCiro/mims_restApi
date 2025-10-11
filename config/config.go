package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             int
	Debug            bool
	GmailEmail       string
	GmailAppPassword string
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

	return &Config{
		Port:             port,
		Debug:            debug,
		GmailEmail:       getEnv("GMAIL_EMAIL", ""),
		GmailAppPassword: getEnv("GMAIL_APP_PASSWORD", ""),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
