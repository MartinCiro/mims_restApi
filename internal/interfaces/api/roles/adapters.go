package roles

import "api_go/internal/core/roles"

// ToRolData convierte DTO de request a entidad del core
func ToRolData(dto CreateRolRequest) roles.RolData {
	return roles.RolData{
		Nombre:      dto.Nombre,
		Descripcion: dto.Descripcion,
		Permisos:    dto.Permisos,
	}
}

// ToRolDataUpdate convierte DTO de update a entidad del core
func ToRolDataUpdate(dto UpdateRolRequest) roles.RolDataUpdate {
	return roles.RolDataUpdate{
		ID:          dto.ID,
		Nombre:      dto.Nombre,
		Descripcion: dto.Descripcion,
		Permisos:    dto.Permisos,
	}
}

// ToRolDataXid convierte DTO de ID a entidad del core
func ToRolDataXid(dto GetRolRequest) roles.RolDataXid {
	return roles.RolDataXid{
		ID: dto.ID,
	}
}

// FromRol convierte entidad del core a DTO de response
func FromRol(rol *roles.Rol) RolResponse {
	return RolResponse{
		ID:          rol.ID,
		Nombre:      rol.Nombre,
		Descripcion: rol.Descripcion,
		Permisos:    rol.Permisos,
	}
}

// FromRoles convierte slice de entidades a DTOs de response
func FromRoles(roles []roles.Rol) []RolResponse {
	responses := make([]RolResponse, len(roles))
	for i, rol := range roles {
		responses[i] = FromRol(&rol)
	}
	return responses
}
