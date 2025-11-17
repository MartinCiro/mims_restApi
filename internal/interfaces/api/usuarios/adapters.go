package usuarios

import "api_go/internal/core/usuarios"

// ToUsuarioData convierte DTO de request a entidad del core
func ToUsuarioData(dto CreateUsuarioRequest) usuarios.UsuarioData {
	return usuarios.UsuarioData{
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
		RolID:    dto.RolID,
		EstadoID: dto.EstadoID,
	}
}

// ToUsuarioDataUpdate convierte DTO de update a entidad del core
func ToUsuarioDataUpdate(dto UpdateUsuarioRequest) usuarios.UsuarioDataUpdate {
	return usuarios.UsuarioDataUpdate{
		Documento: dto.documento,
		Username:  dto.Username,
		Email:     dto.Email,
		RolID:     dto.RolID,
		EstadoID:  dto.EstadoID,
	}
}

// ToUsuarioDataXid convierte DTO de ID a entidad del core
func ToUsuarioDataXid(dto GetUsuarioRequest) usuarios.UsuarioDataXid {
	return usuarios.UsuarioDataXid{
		Documento: dto.documento,
	}
}

// ToCambiarPasswordData convierte DTO de cambio de password a entidad del core
func ToCambiarPasswordData(dto CambiarPasswordRequest) usuarios.CambiarPasswordData {
	return usuarios.CambiarPasswordData{
		ID:             dto.ID,
		PasswordActual: dto.PasswordActual,
		NuevoPassword:  dto.NuevoPassword,
	}
}

// FromUsuario convierte entidad del core a DTO de response
func FromUsuario(usuario *usuarios.Usuario) UsuarioResponse {
	return UsuarioResponse{
		Usuario:   usuario.NombreCompleto,
		Username:  usuario.Username,
		Email:     usuario.Email,
		Rol:       usuario.RolNombre,
		Estado:    usuario.EstadoNombre,
		CreadoAt:  usuario.FechaRegistro,
		UpdatedAt: usuario.FechaActualizacion,
	}
}

// FromUsuarios convierte slice de entidades a DTOs de response
func FromUsuarios(usuarios []usuarios.Usuario) []UsuarioResponse {
	responses := make([]UsuarioResponse, len(usuarios))
	for i, usuario := range usuarios {
		responses[i] = FromUsuario(&usuario)
	}
	return responses
}

// FromUsuarioConRelaciones convierte entidad con relaciones a DTO de response
func FromUsuarioConRelaciones(usuario *usuarios.UsuarioConRelaciones) UsuarioConRelacionesResponse {
	return UsuarioConRelacionesResponse{
		ID:        usuario.Documento,
		Username:  usuario.Username,
		Email:     usuario.Email,
		RolID:     usuario.RolID,
		EstadoID:  usuario.EstadoID,
		RolNombre: usuario.RolNombre,
		Estado:    usuario.Estado,
	}
}

// FromUsuariosConRelaciones convierte slice de entidades con relaciones a DTOs
func FromUsuariosConRelaciones(usuarios []usuarios.UsuarioConRelaciones) []UsuarioConRelacionesResponse {
	responses := make([]UsuarioConRelacionesResponse, len(usuarios))
	for i, usuario := range usuarios {
		responses[i] = FromUsuarioConRelaciones(&usuario)
	}
	return responses
}
