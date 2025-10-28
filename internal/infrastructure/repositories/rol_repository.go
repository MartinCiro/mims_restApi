// internal/infrastructure/repositories/rol_repository.go
package repositories

import (
	"api_go/internal/infrastructure/database"
	"context"
	"fmt"

	"gorm.io/gorm"
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
	fmt.Printf("🔍 Buscando rol: %d\n", id)

	var rolDB struct {
		ID     int    `gorm:"column:id"`
		Nombre string `gorm:"column:nombre_rol"`
	}

	err := r.dbManager.FindUnique(ctx, "rol", &rolDB, map[string]interface{}{
		"id": id,
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("❌ Rol no encontrado: %d\n", id)
			return nil, nil
		}
		fmt.Printf("❌ Error buscando rol: %v\n", err)
		return nil, fmt.Errorf("error buscando rol: %v", err)
	}

	return &Rol{
		ID:     rolDB.ID,
		Nombre: rolDB.Nombre,
	}, nil
}

type Rol struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre_rol"`
}
