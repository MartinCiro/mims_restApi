package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCache(addr, password string, db int) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// TTL por defecto: 1 hora (3600 segundos)
	defaultTTL := 3600 * time.Second

	return &Cache{
		client: client,
		ttl:    defaultTTL,
	}
}

// NewCacheWithConfig crea una instancia con configuración personalizada
func NewCacheWithConfig(addr, password string, db int, ttlSeconds int) *Cache {
	cache := NewCache(addr, password, db)
	cache.ttl = time.Duration(ttlSeconds) * time.Second
	return cache
}

// Set almacena un valor en Redis con TTL
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	// Usar TTL proporcionado o el por defecto
	expiration := c.ttl
	if len(ttl) > 0 {
		expiration = ttl[0]
	}

	// Serializar valor a JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error serializando valor para Redis: %v", err)
	}

	err = c.client.Set(ctx, key, jsonValue, expiration).Err()
	if err != nil {
		return fmt.Errorf("error guardando en Redis: %v", err)
	}

	return nil
}

// Get obtiene un valor de Redis y lo deserializa
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key no existe, no es error
	}
	if err != nil {
		return "", fmt.Errorf("error obteniendo de Redis: %v", err)
	}

	return val, nil
}

// GetJSON obtiene un valor y lo deserializa a una estructura específica
func (c *Cache) GetJSON(ctx context.Context, key string, target interface{}) error {
	val, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	if val == "" {
		return nil // Key no existe
	}

	if err := json.Unmarshal([]byte(val), target); err != nil {
		return fmt.Errorf("error deserializando JSON de Redis: %v", err)
	}

	return nil
}

// Delete elimina una clave de Redis
func (c *Cache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error eliminando de Redis: %v", err)
	}
	return nil
}

// DeletePattern elimina todas las claves que coincidan con un patrón
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("error buscando claves con patrón: %v", err)
	}

	if len(keys) > 0 {
		err = c.client.Del(ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("error eliminando claves con patrón: %v", err)
		}
	}

	return nil
}

// Exists verifica si una clave existe
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("error verificando existencia: %v", err)
	}
	return count > 0, nil
}

// Increment incrementa un valor numérico
func (c *Cache) Increment(ctx context.Context, key string, value int64, ttl ...time.Duration) (int64, error) {
	result, err := c.client.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, fmt.Errorf("error incrementando valor: %v", err)
	}

	// Si se proporciona TTL, establecer expiración
	if len(ttl) > 0 {
		c.client.Expire(ctx, key, ttl[0])
	}

	return result, nil
}

// HealthCheck verifica la conexión con Redis
func (c *Cache) HealthCheck(ctx context.Context) error {
	_, err := c.client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("error en health check de Redis: %v", err)
	}
	return nil
}

// Close cierra la conexión con Redis
func (c *Cache) Close() error {
	return c.client.Close()
}

// GetClient retorna el cliente de Redis (para operaciones avanzadas)
func (c *Cache) GetClient() *redis.Client {
	return c.client
}
