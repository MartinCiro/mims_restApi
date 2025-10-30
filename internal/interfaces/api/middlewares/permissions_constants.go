package middlewares

// Permisos para Estados
var (
	PermissionEstadosListar   = []string{"estados:listar", "estados:leer"}
	PermissionEstadosCrear    = []string{"estados:crear", "estados:escribir"}
	PermissionEstadosEditar   = []string{"estados:editar", "estados:escribir"}
	PermissionEstadosEliminar = []string{"estados:eliminar", "estados:escribir"}
	PermissionEstadosVer      = []string{"estados:ver", "estados:leer"}
)

// Permisos para Usuarios
var (
	PermissionUsuariosListar = []string{"usuarios:listar", "usuarios:leer"}
	PermissionUsuariosCrear  = []string{"usuarios:crear", "usuarios:escribir"}
	PermissionUsuariosEditar = []string{"usuarios:editar", "usuarios:escribir"}
)

// Permisos administrativos
var (
	PermissionAdmin = []string{"admin", "administrador"}
)
