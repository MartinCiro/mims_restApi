package estados

import "context"

type EstadoService struct {
	estadoPort EstadosPort
}

func NewEstadoService(estadoPort EstadosPort) *EstadoService {
	return &EstadoService{
		estadoPort: estadoPort,
	}
}

func (s *EstadoService) ObtenerEstados(ctx context.Context) ([]Estado, error) {
	return s.estadoPort.ObtenerEstados(ctx)
}

func (s *EstadoService) CrearEstado(ctx context.Context, estadoData EstadoData) (*Estado, error) {
	return s.estadoPort.CrearEstados(ctx, estadoData)
}

func (s *EstadoService) ObtenerEstadoXid(ctx context.Context, estadoData EstadoDataXid) (*Estado, error) {
	return s.estadoPort.ObtenerEstadosXid(ctx, estadoData)
}

func (s *EstadoService) UpEstado(ctx context.Context, estadoData EstadoDataUpdate) (*Estado, error) {
	return s.estadoPort.ActualizaEstado(ctx, estadoData)
}

func (s *EstadoService) DelEstado(ctx context.Context, estadoData EstadoDataXid) error {
	return s.estadoPort.DelEstado(ctx, estadoData)
}
