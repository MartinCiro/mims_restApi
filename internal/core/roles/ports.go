package roles

import "context"

// RolesPort define el contrato para la gestión de roles
type RolesPort interface {
	CrearRol(ctx context.Context, rolData RolData) (*Rol, error)
	ObtenerRoles(ctx context.Context) ([]Rol, error)
	ObtenerRolXid(ctx context.Context, rolData RolDataXid) (*Rol, error)
	ActualizarRol(ctx context.Context, rolData RolDataUpdate) (*Rol, error)
	EliminarRol(ctx context.Context, rolData RolDataXid) error
	ObtenerTodosLosPermisos(ctx context.Context) ([]string, error)
}
