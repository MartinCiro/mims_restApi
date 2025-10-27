package estados

import "context"

// EstadosPort define el contrato para la gestión de estados
type EstadosPort interface {
	CrearEstados(ctx context.Context, estadoData EstadoData) (*Estado, error)
	ObtenerEstados(ctx context.Context) ([]Estado, error)
	ObtenerEstadosXid(ctx context.Context, estadoData EstadoDataXid) (*Estado, error)
	DelEstado(ctx context.Context, estadoData EstadoDataXid) error
	ActualizaEstado(ctx context.Context, estadoData EstadoDataUpdate) (*Estado, error)
}
