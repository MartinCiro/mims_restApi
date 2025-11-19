package permisos

import core "api_go/internal/core/permisos"

// ToPermisoData convierte DTO de request a entidad del core
func ToPermisoData(dto CreatePermisoRequest) core.PermisoData {
	return core.PermisoData{
		PermisoBase: core.PermisoBase{
			Nombre:      dto.Nombre,
			Descripcion: &dto.Descripcion,
		},
	}
}

// ToPermisoDataUpdate convierte DTO de update a entidad del core
func ToPermisoDataUpdate(dto UpdatePermisoRequest) core.PermisoDataUpdate {
	return core.PermisoDataUpdate{
		ID: dto.ID,
		PermisoBase: core.PermisoBase{
			Nombre:      dto.Nombre,
			Descripcion: &dto.Descripcion,
		},
	}
}

// ToPermisoDataXid convierte DTO de ID a entidad del core
func ToPermisoDataXid(dto GetPermisoRequest) core.PermisoDataXid {
	return core.PermisoDataXid{
		ID: dto.ID,
	}
}

// FromPermiso convierte entidad del core a DTO de response
func FromPermiso(permiso *core.Permiso) PermisoResponse {
	return PermisoResponse{
		ID: permiso.ID,
		PermisoBase: PermisoBase{
			Nombre:      permiso.Nombre,
			Descripcion: *permiso.Descripcion,
		},
	}
}

// FromPermisos convierte slice de entidades a DTOs de response
func FromPermisos(permisos []core.Permiso) []PermisoResponse {
	responses := make([]PermisoResponse, len(permisos))
	for i, permiso := range permisos {
		responses[i] = FromPermiso(&permiso)
	}
	return responses
}
