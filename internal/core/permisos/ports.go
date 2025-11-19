package permisos

import "context"

// PermisosPort define el contrato para la gestión de permisos
type PermisosPort interface {
	CrearPermisos(ctx context.Context, permisoData PermisoData) (*Permiso, error)
	ObtenerPermisos(ctx context.Context) ([]Permiso, error)
	ObtenerPermisosXid(ctx context.Context, permisoData PermisoDataXid) (*Permiso, error)
	DelPermiso(ctx context.Context, permisoData PermisoDataXid) error
	ActualizaPermiso(ctx context.Context, permisoData PermisoDataUpdate) (*Permiso, error)
}
