package estados

// CreateEstadoRequest DTO para crear estado (desde API)
type CreateEstadoRequest struct {
	Nombre      string  `json:"nombre" binding:"required"`
	Descripcion *string `json:"descripcion,omitempty"`
}

// UpdateEstadoRequest DTO para actualizar estado (desde API)
type UpdateEstadoRequest struct {
	ID          int     `json:"id" binding:"required"`
	Nombre      string  `json:"nombre" binding:"required"`
	Descripcion *string `json:"descripcion,omitempty"`
}

// GetEstadoRequest DTO para obtener estado por ID (desde API)
type GetEstadoRequest struct {
	ID int `json:"id" binding:"required"`
}

// EstadoResponse DTO para respuesta de estado (hacia API)
type EstadoResponse struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion *string `json:"descripcion,omitempty"`
}
