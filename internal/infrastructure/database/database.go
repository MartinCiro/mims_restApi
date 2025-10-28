package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	ServerDB   string
	PortDB     int
	UserDB     string
	PasswordDB string
	Database   string
	SSLMode    string
}

// GetConnection crea y retorna una conexión a la base de datos usando GORM
func GetConnection(config DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		config.ServerDB,
		config.UserDB,
		config.PasswordDB,
		config.Database,
		config.PortDB,
		config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Ajustar según entorno
	})
	if err != nil {
		return nil, fmt.Errorf("error conectando a la base de datos: %v\nCredenciales:\n- Host: %s\n- User: %s\n- Database: %s\n- Port: %d\n- SSLMode: %s",
			err, config.ServerDB, config.UserDB, config.Database, config.PortDB, config.SSLMode)
	}

	// Configurar conexión pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo sql.DB: %v", err)
	}

	// Configurar pool de conexiones
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * 60) // 5 minutos

	return db, nil
}

// HealthCheck verifica que la conexión a la BD esté funcionando
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("error obteniendo sql.DB: %v", err)
	}

	return sqlDB.Ping()
}

// CloseConnection cierra la conexión a la base de datos
func CloseConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("error obteniendo sql.DB: %v", err)
	}

	return sqlDB.Close()
}
