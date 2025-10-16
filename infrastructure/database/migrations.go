package database

import (
	"log"
	"scrapper_go_email/infrastructure/database/models"
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
