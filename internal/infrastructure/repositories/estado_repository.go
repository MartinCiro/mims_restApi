package repositories

import (
	"api_go/internal/core/auth"
	"api_go/internal/infrastructure/database"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Asegurar que EstadoRepository implemente el port
var _ auth.EstadoRepositoryPort = (*EstadoRepository)(nil)

type EstadoRepository struct {
	dbManager *database.DBManager
}

func NewEstadoRepository(dbManager *database.DBManager) *EstadoRepository {
	return &EstadoRepository{
		dbManager: dbManager,
	}
}

// FindEstadoIDByName busca un estado por nombre y retorna su ID
func (r *EstadoRepository) FindEstadoIDByName(ctx context.Context, nombre string) (int, error) {
	fmt.Printf("🔍 Buscando estado por nombre: %s\n", nombre)

	var estadoDB struct {
		ID     int    `gorm:"column:id"`
		Nombre string `gorm:"column:nombre_estado"`
	}

	err := r.dbManager.FindUniqueByField(ctx, "estados", &estadoDB, "nombre_estado", nombre)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("❌ Estado no encontrado: %s\n", nombre)
			return 0, nil
		}
		fmt.Printf("❌ Error buscando estado por nombre: %v\n", err)
		return 0, fmt.Errorf("error buscando estado por nombre: %v", err)
	}

	fmt.Printf("✅ Estado encontrado: ID=%d, Nombre=%s\n", estadoDB.ID, estadoDB.Nombre)
	return estadoDB.ID, nil
}

// FindByID busca un estado por ID
func (r *EstadoRepository) FindByID(ctx context.Context, id int) (*Estado, error) {
	fmt.Printf("🔍 Buscando estado por ID: %d\n", id)

	var estadoDB struct {
		ID     int    `gorm:"column:id"`
		Nombre string `gorm:"column:nombre"`
	}

	err := r.dbManager.FindUniqueByField(ctx, "estados", &estadoDB, "id", id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("❌ Estado no encontrado: %d\n", id)
			return nil, nil
		}
		fmt.Printf("❌ Error buscando estado por ID: %v\n", err)
		return nil, fmt.Errorf("error buscando estado por ID: %v", err)
	}

	return &Estado{
		ID:     estadoDB.ID,
		Nombre: estadoDB.Nombre,
	}, nil
}

// Estado estructura simple para el repositorio
type Estado struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}
