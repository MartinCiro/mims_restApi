package redis

import (
	"api_go/config"
	"context"
	"log"
)

// InitializeCache inicializa y retorna el servicio de Redis con configuración
func InitializeCache(cfg *config.Config) *Cache {
	return InitializeCacheWithOptions(cfg, "")
}

// InitializeCacheWithOptions inicializa Redis con configuración y opciones adicionales
func InitializeCacheWithOptions(cfg *config.Config, customURL string) *Cache {
	var cache *Cache

	if customURL != "" {
		// Usar URL personalizada (override)
		tempCfg := &config.Config{
			RedisURL: customURL,
			RedisTTL: cfg.RedisTTL, // Mantener el TTL de la configuración principal
		}
		cache = NewCache(tempCfg)
	} else {
		// Usar configuración normal
		cache = NewCache(cfg)
	}

	// Verificar conexión
	ctx := context.Background()
	if err := cache.HealthCheck(ctx); err != nil {
		log.Printf("⚠️  No se pudo conectar a Redis: %v", err)
		log.Printf("⚠️  La aplicación funcionará sin cache de Redis")
		// No hacemos fatal para que la app pueda funcionar sin Redis
	} else {
		log.Printf("✅ Conexión a Redis establecida correctamente (TTL: %d segundos)", cfg.RedisTTL)
	}

	return cache
}
