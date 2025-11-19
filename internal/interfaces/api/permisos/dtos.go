package permisos

// Entidades del dominio
type PermisoBase struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion,omitempty"`
}

type Permiso struct {
	ID int `json:"id"`
	PermisoBase
}

type PermisoData struct {
	PermisoBase
}

type PermisoDataUpdate struct {
	ID int `json:"id"`
	PermisoBase
}

type PermisoDataXid struct {
	ID int `json:"id"`
}

// DTOs para la API
type CreatePermisoRequest struct {
	PermisoBase
}

type UpdatePermisoRequest struct {
	ID int `json:"id" validate:"required,min=1"`
	PermisoBase
}

type GetPermisoRequest struct {
	ID int `json:"id" validate:"required,min=1"`
}

type PermisoResponse struct {
	ID int `json:"id"`
	PermisoBase
}
