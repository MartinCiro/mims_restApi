package repositories

import (
	"api_go/internal/infrastructure/database"
	"context"
	"fmt"
)

type PermisoRepository struct {
	dbManager *database.DBManager
}

func NewPermisoRepository(dbManager *database.DBManager) *PermisoRepository {
	return &PermisoRepository{
		dbManager: dbManager,
	}
}

func (r *PermisoRepository) FindByRolID(ctx context.Context, rolID int) ([]string, error) {
	fmt.Printf("🔍 Buscando permisos para rol: %d\n", rolID)

	var permisosDB []struct {
		Nombre string `gorm:"column:nombre_permiso"`
	}

	// Usar join explícito
	query := `
		SELECT permiso.nombre_permiso
		FROM permiso 
		INNER JOIN rol_x_permiso ON permiso.id = rol_x_permiso.id_permiso 
		WHERE rol_x_permiso.id_rol = ?
	`

	err := r.dbManager.FindWithJoin(ctx, &permisosDB, query, rolID)
	if err != nil {
		fmt.Printf("❌ Error buscando permisos: %v\n", err)
		return nil, fmt.Errorf("error buscando permisos: %v", err)
	}

	permisos := make([]string, len(permisosDB))
	for i, p := range permisosDB {
		permisos[i] = p.Nombre
	}

	fmt.Printf("✅ Permisos encontrados: %d\n", len(permisos))
	return permisos, nil
}
