package estados

import "api_go/internal/core/estados"

// ToEstadoData convierte DTO de request a entidad del core
func ToEstadoData(dto CreateEstadoRequest) estados.EstadoData {
	return estados.EstadoData{
		Nombre:      dto.Nombre,
		Descripcion: dto.Descripcion,
	}
}

// ToEstadoDataUpdate convierte DTO de update a entidad del core
func ToEstadoDataUpdate(dto UpdateEstadoRequest) estados.EstadoDataUpdate {
	return estados.EstadoDataUpdate{
		ID:          dto.ID,
		Nombre:      dto.Nombre,
		Descripcion: dto.Descripcion,
	}
}

// ToEstadoDataXid convierte DTO de ID a entidad del core
func ToEstadoDataXid(dto GetEstadoRequest) estados.EstadoDataXid {
	return estados.EstadoDataXid{
		ID: dto.ID,
	}
}

// FromEstado convierte entidad del core a DTO de response
func FromEstado(estado *estados.Estado) EstadoResponse {
	return EstadoResponse{
		ID:          estado.ID,
		Nombre:      estado.Nombre,
		Descripcion: estado.Descripcion,
	}
}

// FromEstados convierte slice de entidades a DTOs de response
func FromEstados(estados []estados.Estado) []EstadoResponse {
	responses := make([]EstadoResponse, len(estados))
	for i, estado := range estados {
		responses[i] = FromEstado(&estado)
	}
	return responses
}
