package roles

// CreateRolRequest DTO para crear rol (desde API)
type CreateRolRequest struct {
	Nombre      string   `json:"nombre" binding:"required"`
	Descripcion string   `json:"descripcion" binding:"required"`
	Permisos    []string `json:"permisos" binding:"required"`
}

// UpdateRolRequest DTO para actualizar rol (desde API)
type UpdateRolRequest struct {
	ID          int      `json:"id" binding:"required"`
	Nombre      string   `json:"nombre" binding:"required"`
	Descripcion string   `json:"descripcion" binding:"required"`
	Permisos    []string `json:"permisos" binding:"required"`
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
