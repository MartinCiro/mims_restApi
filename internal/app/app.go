package app

import (
	"context"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"api_go/config"
	"api_go/infrastructure/database/models"
	"api_go/internal/core/auth"
	"api_go/internal/core/estados"
	"api_go/internal/core/login"
	"api_go/internal/infrastructure/database"
	"api_go/internal/infrastructure/jwt"
	"api_go/internal/infrastructure/redis"
	"api_go/internal/infrastructure/repositories"
	"api_go/pkg/logger"
	"api_go/pkg/utils"
)

type App struct {
	// Configuración
	Config *config.Config

	// Infraestructura
	DB         *gorm.DB
	RedisCache *redis.Cache
	JWTService *jwt.JWTService

	// Servicios del Core
	AuthService   *auth.AuthService
	EstadoService *estados.EstadoService

	// Servicios de Utilidad
	PasswordService *utils.PasswordService

	LoginService *login.LoginService
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize(cfg *config.Config) error {
	a.Config = cfg

	// Inicializar infraestructura
	if err := a.initializeInfrastructure(); err != nil {
		return err
	}

	// Inicializar servicios
	if err := a.initializeServices(); err != nil {
		return err
	}

	// Ejecutar migraciones si es necesario
	if err := a.runMigrations(); err != nil {
		logger.Error("❌ Error en migraciones", "error", err)
		// En producción podrías continuar, en desarrollo fallar
		if cfg.IsDevelopment() {
			return fmt.Errorf("error en migraciones: %v", err)
		}
	}

	/* log.Println("✅ Todos los servicios inicializados correctamente") */
	return nil
}

func (a *App) initializeInfrastructure() error {
	// Inicializar base de datos
	db, err := a.initializeDatabase()
	if err != nil {
		return fmt.Errorf("error inicializando base de datos: %v", err)
	}
	a.DB = db

	// Inicializar Redis
	a.RedisCache = a.initializeRedis()

	// Inicializar JWT Service
	a.JWTService = a.initializeJWTService()

	return nil
}

func (a *App) initializeDatabase() (*gorm.DB, error) {
	// Convertir PortDB a entero
	port, err := strconv.Atoi(a.Config.PortDB)
	if err != nil {
		port = 5432 // Valor por defecto
		logger.Warn("usando puerto por defecto para BD", "port", port)
	}

	dbConfig := database.DBConfig{
		ServerDB:   a.Config.ServerDB,
		PortDB:     port,
		UserDB:     a.Config.UserDB,
		PasswordDB: a.Config.PasswordDB,
		Database:   a.Config.Database,
		SSLMode:    "disable",
	}

	db, err := database.GetConnection(dbConfig)
	if err != nil {
		return nil, err
	}

	// Verificar conexión
	if err := database.HealthCheck(db); err != nil {
		return nil, fmt.Errorf("error en health check de BD: %v", err)
	}

	return db, nil
}

func (a *App) initializeRedis() *redis.Cache {
	cache := redis.InitializeCache(a.Config)
	logger.Info("redis configurado", "url", a.Config.RedisURL, "ttl", a.Config.RedisTTL)
	return cache
}

func (a *App) initializeJWTService() *jwt.JWTService {
	if a.Config.JWTSecret == "" {
		logger.Warn("JWTSecret está vacío en la configuración")
	} else {
		logger.Info("JWTSecret cargado correctamente", "length", len(a.Config.JWTSecret))
	}
	return jwt.NewJWTService(
		a.Config.JWTSecret,
		a.Config.Env,
	)
}

func (a *App) initializeServices() error {
	// Inicializar servicios de utilidad
	a.PasswordService = utils.NewPasswordService(a.Config.JWTSalt)

	// Inicializar DB Manager corregido
	dbManager := database.NewDBManager(a.DB)

	// Inicializar repositorios separados
	usuarioRepo := repositories.NewUsuarioRepository(dbManager)
	rolRepo := repositories.NewRolRepository(dbManager)
	permisoRepo := repositories.NewPermisoRepository(dbManager)

	// AuthAdapter coordina los repositorios
	authAdapter := repositories.NewAuthAdapter(usuarioRepo, rolRepo, permisoRepo)

	// Inicializar servicios del core
	a.AuthService = auth.NewAuthService(
		authAdapter,
		a.RedisCache,
		a.JWTService,
		a.PasswordService,
	)

	a.LoginService = login.NewLoginService(
		authAdapter,
		a.JWTService,
		a.RedisCache,
		a.PasswordService,
		a.Config,
	)

	// Servicio de estados
	estadosAdapter := repositories.NewEstadosAdapter(a.DB, a.RedisCache)
	a.EstadoService = estados.NewEstadoService(estadosAdapter)

	logger.Info("servicios del core inicializados correctamente")
	return nil
}

func (a *App) runMigrations() error {
	logger.Info("🔄 Iniciando migraciones...")

	models := []interface{}{
		&models.Migration{},
		&models.Estado{},
		&models.Rol{},
		&models.Permiso{},
		&models.RolXPermiso{},
		&models.Usuario{},
	}

	for _, model := range models {
		logger.Info("Migrando modelo", "model", fmt.Sprintf("%T", model))
		if err := a.DB.AutoMigrate(model); err != nil {
			logger.Error("Error migrando modelo", "model", fmt.Sprintf("%T", model), "error", err)
			return err
		}
		logger.Info("✅ Modelo migrado", "model", fmt.Sprintf("%T", model))
	}

	logger.Info("✅ Todas las migraciones completadas")
	return nil
}

func (a *App) Shutdown() {
	logger.Info("cerrando servicios de la aplicación")

	// Cerrar base de datos
	if a.DB != nil {
		if sqlDB, err := a.DB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				logger.Warn("error cerrando conexión de BD", "error", err)
			}
		}
	}

	// Cerrar Redis
	if a.RedisCache != nil {
		if err := a.RedisCache.Close(); err != nil {
			logger.Warn("error cerrando conexión de Redis", "error", err)
		}
	}

	logger.Info("todos los servicios cerrados correctamente")
}

// HealthCheck verifica el estado de todos los servicios
func (a *App) HealthCheck() error {
	// Verificar base de datos
	if err := database.HealthCheck(a.DB); err != nil {
		return fmt.Errorf("base de datos no saludable: %v", err)
	}

	// Verificar Redis
	if err := a.RedisCache.HealthCheck(context.Background()); err != nil {
		return fmt.Errorf("redis no saludable: %v", err)
	}

	return nil
}

// GetConfig retorna la configuración (útil para tests)
func (a *App) GetConfig() *config.Config {
	return a.Config
}

// GetDB retorna la instancia de BD (útil para tests y operaciones directas)
func (a *App) GetDB() *gorm.DB {
	return a.DB
}

// GetRedis retorna el cliente de Redis (útil para operaciones avanzadas)
func (a *App) GetRedis() *redis.Cache {
	return a.RedisCache
}
