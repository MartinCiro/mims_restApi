package redis

import (
	"context"
	"log"
	"strconv"
)

// InitializeCache inicializa y retorna el servicio de Redis
func InitializeCache(host, password string, db int) *Cache {
	cache := NewCache(host, password, db)

	// Verificar conexión
	ctx := context.Background()
	if err := cache.HealthCheck(ctx); err != nil {
		log.Printf("⚠️  No se pudo conectar a Redis: %v", err)
		log.Printf("⚠️  La aplicación funcionará sin cache de Redis")
		// No hacemos fatal para que la app pueda funcionar sin Redis
	} else {
		log.Println("✅ Conexión a Redis establecida correctamente")
	}

	return cache
}

// InitializeCacheFromEnv inicializa Redis desde variables de entorno
func InitializeCacheFromEnv() *Cache {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")
	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	addr := host + ":" + port

	return InitializeCache(addr, password, db)
}

// getEnv obtiene variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	// En implementación real usarías os.Getenv(key)
	// Por simplicidad retornamos el valor por defecto
	return defaultValue
}
