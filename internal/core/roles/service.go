package roles

import (
	"context"
	"fmt"
	"strings"
)

type RolService struct {
	rolPort RolesPort
}

func NewRolService(rolPort RolesPort) *RolService {
	return &RolService{
		rolPort: rolPort,
	}
}

func (s *RolService) CrearRol(ctx context.Context, rolData RolData) (*Rol, error) {
	// Validar que el nombre no esté vacío
	if strings.TrimSpace(rolData.Nombre) == "" {
		return nil, fmt.Errorf("status_cod:400, data:El nombre del rol es requerido")
	}

	// Validar que la descripción no esté vacía
	if strings.TrimSpace(rolData.Descripcion) == "" {
		return nil, fmt.Errorf("status_cod:400, data:La descripción del rol es requerida")
	}

	// Validar que haya permisos
	if len(rolData.Permisos) == 0 {
		return nil, fmt.Errorf("status_cod:400, data:Se requiere al menos un permiso")
	}

	return s.rolPort.CrearRol(ctx, rolData)
}

func (s *RolService) ObtenerRoles(ctx context.Context) ([]Rol, error) {
	return s.rolPort.ObtenerRoles(ctx)
}

func (s *RolService) ObtenerRolXid(ctx context.Context, rolData RolDataXid) (*Rol, error) {
	if rolData.ID <= 0 {
		return nil, fmt.Errorf("status_cod:400, data:ID de rol inválido")
	}
	return s.rolPort.ObtenerRolXid(ctx, rolData)
}

func (s *RolService) ActualizarRol(ctx context.Context, rolData RolDataUpdate) (*Rol, error) {
	if rolData.ID <= 0 {
		return nil, fmt.Errorf("status_cod:400, data:ID de rol inválido")
	}

	// Validar que el nombre no esté vacío
	if strings.TrimSpace(rolData.Nombre) == "" {
		return nil, fmt.Errorf("status_cod:400, data:El nombre del rol es requerido")
	}

	// Validar que la descripción no esté vacía
	if strings.TrimSpace(rolData.Descripcion) == "" {
		return nil, fmt.Errorf("status_cod:400, data:La descripción del rol es requerida")
	}

	// Validar que haya permisos
	if len(rolData.Permisos) == 0 {
		return nil, fmt.Errorf("status_cod:400, data:Se requiere al menos un permiso")
	}

	return s.rolPort.ActualizarRol(ctx, rolData)
}

func (s *RolService) EliminarRol(ctx context.Context, rolData RolDataXid) error {
	if rolData.ID <= 0 {
		return fmt.Errorf("status_cod:400, data:ID de rol inválido")
	}
	return s.rolPort.EliminarRol(ctx, rolData)
}

func (s *RolService) ObtenerTodosLosPermisos(ctx context.Context) ([]string, error) {
	return s.rolPort.ObtenerTodosLosPermisos(ctx)
}
