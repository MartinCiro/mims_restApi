package permisos

import "context"

type PermisoService struct {
	permisoPort PermisosPort
}

func NewPermisoService(permisoPort PermisosPort) *PermisoService {
	return &PermisoService{
		permisoPort: permisoPort,
	}
}

func (s *PermisoService) ObtenerPermisos(ctx context.Context) ([]Permiso, error) {
	return s.permisoPort.ObtenerPermisos(ctx)
}

func (s *PermisoService) CrearPermiso(ctx context.Context, permisoData PermisoData) (*Permiso, error) {
	return s.permisoPort.CrearPermisos(ctx, permisoData)
}

func (s *PermisoService) ObtenerPermisoXid(ctx context.Context, permisoData PermisoDataXid) (*Permiso, error) {
	return s.permisoPort.ObtenerPermisosXid(ctx, permisoData)
}

func (s *PermisoService) UpPermiso(ctx context.Context, permisoData PermisoDataUpdate) (*Permiso, error) {
	return s.permisoPort.ActualizaPermiso(ctx, permisoData)
}

func (s *PermisoService) DelPermiso(ctx context.Context, permisoData PermisoDataXid) error {
	return s.permisoPort.DelPermiso(ctx, permisoData)
}
