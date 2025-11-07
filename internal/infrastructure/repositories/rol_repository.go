package repositories

import (
	"api_go/internal/core/auth"
	"api_go/internal/infrastructure/database"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Asegurar que RolRepository implemente el port
var _ auth.RolRepositoryPort = (*RolRepository)(nil)

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

	err := r.dbManager.FindUnique(ctx, "roles", &rolDB, map[string]interface{}{
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

// FindRolIDByName busca un rol por nombre y retorna su ID
func (r *RolRepository) FindRolIDByName(ctx context.Context, nombre string) (int, error) {
	fmt.Printf("🔍 Buscando rol por nombre: %s\n", nombre)

	var rolDB struct {
		ID     int    `gorm:"column:id"`
		Nombre string `gorm:"column:nombre"`
	}

	err := r.dbManager.FindUniqueByField(ctx, "roles", &rolDB, "nombre_rol", nombre)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("❌ Rol no encontrado: %s\n", nombre)
			return 0, nil
		}
		fmt.Printf("❌ Error buscando rol por nombre: %v\n", err)
		return 0, fmt.Errorf("error buscando rol por nombre: %v", err)
	}

	fmt.Printf("✅ Rol encontrado: ID=%d, Nombre=%s\n", rolDB.ID, rolDB.Nombre)
	return rolDB.ID, nil
}

type Rol struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre_rol"`
}
