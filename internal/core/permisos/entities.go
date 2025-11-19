package permisos

// PermisoBase contiene los campos fundamentales de un permiso
type PermisoBase struct {
	Nombre      string  `json:"nombre"`
	Descripcion *string `json:"descripcion,omitempty"`
}

// Permiso representa la entidad completa con ID
type Permiso struct {
	ID int `json:"id"`
	PermisoBase
}

// PermisoData para creación (sin ID)
type PermisoData struct {
	PermisoBase
}

// PermisoDataXid para operaciones por ID
type PermisoDataXid struct {
	ID int `json:"id"`
}

// PermisoDataUpdate para actualización
type PermisoDataUpdate struct {
	ID int `json:"id"`
	PermisoBase
}
