package roles

// Rol representa la entidad de dominio rol
type Rol struct {
	ID          int      `json:"id"`
	Nombre      string   `json:"nombre"`
	Descripcion string   `json:"descripcion"`
	Permisos    []string `json:"permisos,omitempty"`
}

// RolData contiene los datos para crear un rol
type RolData struct {
	Nombre      string   `json:"nombre" binding:"required"`
	Descripcion string   `json:"descripcion" binding:"required"`
	Permisos    []string `json:"permisos" binding:"required"`
}

// RolDataXid contiene el ID para buscar o eliminar un rol
type RolDataXid struct {
	ID int `json:"id" binding:"required"`
}

// RolDataUpdate contiene los datos para actualizar un rol
type RolDataUpdate struct {
	ID          int      `json:"id" binding:"required"`
	Nombre      string   `json:"nombre" binding:"required"`
	Descripcion string   `json:"descripcion" binding:"required"`
	Permisos    []string `json:"permisos" binding:"required"`
}
