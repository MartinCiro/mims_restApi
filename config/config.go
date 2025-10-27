package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Server Config
	Port string
	Env  string

	// Database
	UserDB     string
	PasswordDB string
	ServerDB   string
	Database   string
	PortDB     string

	// Auth
	JWTSecret     string
	JWTSalt       string
	JWTExpireTime int

	// Redis
	RedisTTL int
	RedisURL string
}

// Load carga la configuración desde variables de entorno
func Load() *Config {
	// Cargar .env si existe (usando godotenv si lo prefieres)
	// godotenv.Load()

	return &Config{
		// Server Config
		Port: getEnv("PORT", "3000"),
		Env:  getEnv("ENV", "Production"),

		// Database
		UserDB:     getEnv("USER_DB", ""),
		PasswordDB: getEnv("PASSWORD_DB", ""),
		ServerDB:   getEnv("SERVER_DB", ""),
		Database:   getEnv("DATABASE", ""),
		PortDB:     getEnv("PORT_DB", ""),

		// Auth
		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTSalt:       getEnv("JWT_SALT", "10"),
		JWTExpireTime: getEnvAsInt("JWT_TIEMPO_EXPIRA", 3600),

		// Redis
		RedisTTL: getEnvAsInt("REDIS_TTL", 3600),
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),
	}
}

// getEnv obtiene variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt obtiene variable de entorno como entero
func getEnvAsInt(key string, defaultValue int) int {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}

	return value
}

// IsDevelopment retorna true si el entorno es desarrollo
func (c *Config) IsDevelopment() bool {
	return c.Env == "Dev" || c.Env == "Development"
}

// IsProduction retorna true si el entorno es producción
func (c *Config) IsProduction() bool {
	return c.Env == "Production" || c.Env == "Prod"
}
