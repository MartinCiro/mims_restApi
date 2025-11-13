package roles

// CreateRolRequest DTO para crear rol (desde API)
type CreateRolRequest struct {
	Nombre      string `json:"nombre" binding:"required"`
	Descripcion string `json:"descripcion" binding:"required"`
	Permisos    []int  `json:"permisos" validate:"required,min=1"`
}

type UpdateRolRequest struct {
	ID          int     `json:"id" validate:"required,min=1"`
	Nombre      *string `json:"nombre,omitempty" validate:"omitempty,min=1,max=100"`
	Descripcion *string `json:"descripcion,omitempty" validate:"omitempty"`
	Permisos    []int   `json:"permisos,omitempty" validate:"omitempty,min=1"`
}

// GetRolRequest DTO para obtener rol por ID (desde API)
type GetRolRequest struct {
	ID int `json:"id" binding:"required"`
}

// RolResponse DTO para respuesta de rol (hacia API)
type RolResponse struct {
	ID          int      `json:"id"`
	Nombre      string   `json:"nombre"`
	Descripcion string   `json:"descripcion"`
	Permisos    []string `json:"permisos"`
}
