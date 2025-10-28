// internal/infrastructure/repositories/permiso_repository.go
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

// FindByRolID simula: prisma.permiso.findMany({ where: { rol_x_permiso: { id_rol } } })
func (r *PermisoRepository) FindByRolID(ctx context.Context, rolID int) ([]string, error) {
	var permisosDB []struct {
		Nombre string `gorm:"column:nombre"`
	}

	err := r.dbManager.FindMany(ctx, &permisosDB, map[string]interface{}{
		"rol_x_permiso.id_rol": rolID,
	}, database.WithSelect{Fields: []string{"permiso.nombre"}})

	if err != nil {
		return nil, fmt.Errorf("error buscando permisos: %v", err)
	}

	permisos := make([]string, len(permisosDB))
	for i, p := range permisosDB {
		permisos[i] = p.Nombre
	}

	return permisos, nil
}
