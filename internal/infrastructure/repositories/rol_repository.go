package repositories

import (
	"api_go/infrastructure/database/models"
	"api_go/internal/core/auth"
	"api_go/internal/infrastructure/database"
	"api_go/pkg/logger"
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

// FindByID busca un rol por ID
func (r *RolRepository) FindByID(ctx context.Context, id int) (*models.Rol, error) {
	logger.Info("Buscando rol por ID", "ID", id)

	var rol models.Rol

	err := r.dbManager.GetDB().WithContext(ctx).
		First(&rol, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("❌ Rol no encontrado", "ID", id)
			return nil, nil
		}
		logger.Error("❌ Error buscando rol por ID", "error", err)
		return nil, fmt.Errorf("error buscando rol por ID: %v", err)
	}

	logger.Info("✅ Rol encontrado", "ID", rol.ID, "NombreRol", rol.NombreRol)
	return &rol, nil
}

// FindRolIDByName busca un rol por nombre y retorna su ID
func (r *RolRepository) FindRolIDByName(ctx context.Context, nombre string) (int, error) {
	logger.Info("Buscando rol por nombre", "nombre", nombre)

	var rol models.Rol

	err := r.dbManager.GetDB().WithContext(ctx).
		Model(&models.Rol{}).
		Where("nombre_rol = ?", nombre).
		First(&rol).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("❌ Rol no encontrado", "nombre", nombre)
			return 0, nil
		}
		logger.Error("❌ Error buscando rol por nombre", "error", err)
		return 0, fmt.Errorf("error buscando rol por nombre: %v", err)
	}

	logger.Info("✅ Rol encontrado", "ID", rol.ID, "NombreRol", rol.NombreRol)
	return rol.ID, nil
}

// FindAll busca todos los roles
func (r *RolRepository) FindAll(ctx context.Context) ([]models.Rol, error) {
	logger.Info("Buscando todos los roles")

	var roles []models.Rol

	err := r.dbManager.GetDB().WithContext(ctx).
		Find(&roles).Error

	if err != nil {
		logger.Error("❌ Error buscando todos los roles", "error", err)
		return nil, fmt.Errorf("error buscando todos los roles: %v", err)
	}

	logger.Info("✅ Roles encontrados", "count", len(roles))
	return roles, nil
}

// FindByNombre busca roles por nombre (búsqueda parcial)
func (r *RolRepository) FindByNombre(ctx context.Context, nombre string) ([]models.Rol, error) {
	logger.Info("🔍 Buscando roles por nombre", "nombre", nombre)

	var roles []models.Rol

	err := r.dbManager.GetDB().WithContext(ctx).
		Where("nombre_rol LIKE ?", "%"+nombre+"%").
		Find(&roles).Error

	if err != nil {
		logger.Error("❌ Error buscando roles por nombre", "error", err)
		return nil, fmt.Errorf("error buscando roles por nombre: %v", err)
	}

	logger.Info("✅ Roles encontrados por nombre", "count", len(roles))
	return roles, nil
}

// Create crea un nuevo rol
func (r *RolRepository) Create(ctx context.Context, rol *models.Rol) error {
	logger.Info("Creando nuevo rol", "nombre", rol.NombreRol)

	err := r.dbManager.GetDB().WithContext(ctx).
		Create(rol).Error

	if err != nil {
		logger.Error("❌ Error creando rol", "error", err)
		return fmt.Errorf("error creando rol: %v", err)
	}

	logger.Info("✅ Rol creado", "ID", rol.ID, "NombreRol", rol.NombreRol)
	return nil
}

// Update actualiza un rol existente
func (r *RolRepository) Update(ctx context.Context, rol *models.Rol) error {
	logger.Info("Actualizando rol", "ID", rol.ID)

	err := r.dbManager.GetDB().WithContext(ctx).
		Save(rol).Error

	if err != nil {
		logger.Error("❌ Error actualizando rol", "error", err)
		return fmt.Errorf("error actualizando rol: %v", err)
	}

	logger.Info("✅ Rol actualizado", "ID", rol.ID)
	return nil
}

// Delete elimina un rol por ID
func (r *RolRepository) Delete(ctx context.Context, id int) error {
	logger.Info("Eliminando rol", "ID", id)

	result := r.dbManager.GetDB().WithContext(ctx).
		Delete(&models.Rol{}, id)

	if result.Error != nil {
		logger.Error("❌ Error eliminando rol", "error", result.Error)
		return fmt.Errorf("error eliminando rol: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		logger.Error("❌ Rol no encontrado para eliminar", "ID", id)
		return fmt.Errorf("rol no encontrado")
	}

	logger.Info("✅ Rol eliminado", "ID", id)
	return nil
}

// GetRolesConPermisos obtiene roles con sus permisos asociados
func (r *RolRepository) GetRolesConPermisos(ctx context.Context) ([]models.Rol, error) {
	logger.Info("Buscando roles con permisos")

	var roles []models.Rol

	err := r.dbManager.GetDB().WithContext(ctx).
		Preload("Permisos").
		Find(&roles).Error

	if err != nil {
		logger.Error("❌ Error buscando roles con permisos", "error", err)
		return nil, fmt.Errorf("error buscando roles con permisos: %v", err)
	}

	logger.Info("✅ Roles con permisos encontrados", "count", len(roles))
	return roles, nil
}
