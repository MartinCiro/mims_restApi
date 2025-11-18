package usuarios

import (
	"time"
)

// Usuario representa la entidad de dominio usuario
type Usuario struct {
	Documento          string    `json:"documento"`
	Nombres            string    `json:"nombres"`
	Apellido           string    `json:"apellido"`
	NombreCompleto     string    `json:"usuario"` // nombre + apellido
	Email              string    `json:"email"`
	Username           string    `json:"username"`
	RolID              int       `json:"rol_id"`
	EstadoID           int       `json:"estado_id"`
	RolNombre          string    `json:"rol"`
	EstadoNombre       string    `json:"estado"`
	FechaRegistro      time.Time `json:"fecha_registro"`
	FechaActualizacion time.Time `json:"fecha_actualizacion"`
}

// UsuarioData contiene los datos para crear un nuevo usuario
type UsuarioData struct {
	Documento       string     `json:"documento" binding:"required"`
	Nombres         string     `json:"nombres" binding:"required"`
	Apellido        string     `json:"apellido" binding:"required"`
	Email           string     `json:"email" binding:"required,email"`
	Username        string     `json:"username" binding:"required"`
	Password        string     `json:"password" binding:"required,min=6"`
	RolID           *int       `json:"rol_id,omitempty"`
	EstadoID        *int       `json:"estado_id,omitempty"`
	InfoPerfil      *string    `json:"info_perfil,omitempty"`
	NumContacto     *string    `json:"num_contacto,omitempty"`
	FechaNacimiento *time.Time `json:"fecha_nacimiento,omitempty"`
}

// UsuarioDataXid contiene el ID para buscar o eliminar un usuario
type UsuarioDataXid struct {
	Documento string `json:"documento" binding:"required"`
}

// UsuarioDataUpdate contiene los datos para actualizar un usuario
type UsuarioDataUpdate struct {
	Documento       string     `json:"documento" binding:"required"`
	Username        *string    `json:"username,omitempty"`
	Email           *string    `json:"email,omitempty"`
	Nombres         *string    `json:"nombres,omitempty"`
	Apellido        *string    `json:"apellido,omitempty"`
	InfoPerfil      *string    `json:"info_perfil,omitempty"`
	NumContacto     *string    `json:"num_contacto,omitempty"`
	FechaNacimiento *time.Time `json:"fecha_nacimiento,omitempty"`
	RolID           *int       `json:"rol_id,omitempty"`
	EstadoID        *int       `json:"estado_id,omitempty"`
}

// CambiarPasswordData contiene los datos para cambiar contraseña
type CambiarPasswordData struct {
	ID             int    `json:"id"`
	PasswordActual string `json:"password_actual"`
	NuevoPassword  string `json:"nuevo_password"`
}

// UsuarioConRelaciones representa un usuario con información de relaciones
type UsuarioConRelaciones struct {
	Documento       string     `json:"documento"`
	Nombres         string     `json:"nombres"`
	Apellido        string     `json:"apellido"`
	Email           string     `json:"email"`
	Username        string     `json:"username"`
	RolID           int        `json:"rol_id"`
	EstadoID        int        `json:"estado_id"`
	InfoPerfil      *string    `json:"info_perfil,omitempty"`
	NumContacto     *string    `json:"num_contacto,omitempty"`
	FechaNacimiento *time.Time `json:"fecha_nacimiento,omitempty"`
	FechaRegistro   time.Time  `json:"fecha_registro"`
	RolNombre       string     `json:"rol_nombre"`
	Estado          string     `json:"estado"`
}
