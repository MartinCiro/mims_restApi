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
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Permisos    []int  `json:"permisos"`
}

// RolDataXid contiene el ID para buscar o eliminar un rol
type RolDataXid struct {
	ID int `json:"id"`
}

// RolDataUpdate contiene los datos para actualizar un rol
type RolDataUpdate struct {
	ID          int     `json:"id"`
	Nombre      *string `json:"nombre,omitempty"`
	Descripcion *string `json:"descripcion,omitempty"`
	Permisos    []int   `json:"permisos,omitempty"`
}
