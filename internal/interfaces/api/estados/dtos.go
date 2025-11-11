package estados

type CreateEstadoRequest struct {
	Nombre      string  `json:"nombre" validate:"required,min=1,max=100"`
	Descripcion *string `json:"descripcion,omitempty" validate:"omitempty,max=255"`
}

type UpdateEstadoRequest struct {
	ID          int     `json:"id" validate:"required,min=1"`
	Nombre      string  `json:"nombre" validate:"required,min=1,max=100"`
	Descripcion *string `json:"descripcion,omitempty" validate:"omitempty,max=255"`
}

type GetEstadoRequest struct {
	ID int `json:"id" validate:"required,min=1"`
}

type EstadoResponse struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion *string `json:"descripcion,omitempty"`
}
