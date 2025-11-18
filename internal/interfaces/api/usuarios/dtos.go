package usuarios

import (
	"time"
)

type CreateUsuarioRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	RolID    *int   `json:"rol_id,omitempty"`
	EstadoID *int   `json:"estado_id,omitempty"`
}

type GetUsuarioRequest struct {
	documento string `json:"documento"`
}

// UpdateUsuarioRequest - DTO para actualizar usuario
type UpdateUsuarioRequest struct {
	Documento       string     `json:"-"` // No viene del JSON, se asigna desde la URL
	Username        *string    `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email           *string    `json:"email,omitempty" validate:"omitempty,email"`
	Nombres         *string    `json:"nombres,omitempty" validate:"omitempty,min=2"`
	Apellido        *string    `json:"apellido,omitempty" validate:"omitempty,min=2"`
	InfoPerfil      *string    `json:"info_perfil,omitempty"`
	NumContacto     *string    `json:"num_contacto,omitempty"`
	FechaNacimiento *time.Time `json:"fecha_nacimiento,omitempty"`
	RolID           *int       `json:"rol_id,omitempty"`
	EstadoID        *int       `json:"estado_id,omitempty"`
}

// CambiarPasswordRequest - DTO para cambiar contraseña
type CambiarPasswordRequest struct {
	ID             int    `json:"-"`
	PasswordActual string `json:"password_actual" validate:"required"`
	NuevoPassword  string `json:"nuevo_password" validate:"required,min=6"`
}

// UsuarioResponse - DTO para respuesta de usuario
type UsuarioResponse struct {
	Documento string    `json:"documento"`
	Usuario   string    `json:"usuario"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Rol       string    `json:"rol"`
	Estado    string    `json:"estado"`
	CreadoAt  time.Time `json:"creado_at"`
}

type UsuarioConRelacionesResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	RolID     int       `json:"rol_id"`
	EstadoID  int       `json:"estado_id"`
	RolNombre string    `json:"rol_nombre"`
	Estado    string    `json:"estado"`
	CreatedAt time.Time `json:"created_at"`
}
