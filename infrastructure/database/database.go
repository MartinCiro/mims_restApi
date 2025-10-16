package database

import (
	"log"
	"scrapper_go_email/config"

	"github.com/glebarez/sqlite" // Driver puro en Go
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	var err error

	// Configurar logger de GORM
	gormConfig := &gorm.Config{}
	if cfg.DatabaseLogMode {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	// Conectar a SQLite con driver puro en Go
	DB, err = gorm.Open(sqlite.Open(cfg.DatabaseURL), gormConfig)
	if err != nil {
		return err
	}

	log.Println("✅ Base de datos conectada:", cfg.DatabaseURL)
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
