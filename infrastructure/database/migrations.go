package database

import (
	"api_go/infrastructure/database/models"
	"log"
)

func RunMigrations() error {
	err := DB.AutoMigrate(
		&models.Session{},
		// Agregar otros modelos aquí
	)

	if err != nil {
		return err
	}

	log.Println("✅ Migraciones ejecutadas correctamente")
	return nil
}
