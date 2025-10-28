// internal/infrastructure/repositories/rol_repository.go
package repositories

import (
	"api_go/internal/infrastructure/database"
	"context"
	"fmt"
)

type RolRepository struct {
	dbManager *database.DBManager
}

func NewRolRepository(dbManager *database.DBManager) *RolRepository {
	return &RolRepository{
		dbManager: dbManager,
	}
}

// FindByID simula: prisma.rol.findUnique({ where: { id } })
func (r *RolRepository) FindByID(ctx context.Context, id int) (*Rol, error) {
	var rolDB struct {
		ID     int    `gorm:"column:id"`
		Nombre string `gorm:"column:nombre"`
	}

	err := r.dbManager.FindUnique(ctx, &rolDB, map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return nil, fmt.Errorf("error buscando rol: %v", err)
	}

	return &Rol{
		ID:     rolDB.ID,
		Nombre: rolDB.Nombre,
	}, nil
}

type Rol struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}
