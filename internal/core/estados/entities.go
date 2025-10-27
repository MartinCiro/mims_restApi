package estados

// Estado representa la entidad de dominio estado
type Estado struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion *string `json:"descripcion,omitempty"`
}

// EstadoData contiene los datos para crear un estado
type EstadoData struct {
	Nombre      string  `json:"nombre"`
	Descripcion *string `json:"descripcion,omitempty"`
}

// EstadoDataXid contiene el ID para buscar o eliminar un estado
type EstadoDataXid struct {
	ID int `json:"id"`
}

// EstadoDataUpdate contiene los datos para actualizar un estado
type EstadoDataUpdate struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion *string `json:"descripcion,omitempty"`
}
