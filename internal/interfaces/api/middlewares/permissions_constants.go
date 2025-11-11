package middlewares

// Permisos para Estados
var (
	PermissionEstadosListar   = []string{"estados:listar", "estados:leer"}
	PermissionEstadosCrear    = []string{"estados:crear", "estados:escribir"}
	PermissionEstadosEditar   = []string{"estados:editar", "estados:escribir"}
	PermissionEstadosEliminar = []string{"estados:eliminar", "estados:escribir"}
	PermissionEstadosVer      = []string{"estados:ver", "estados:leer"}
)

// Permisos para Roles
var (
	PermissionRolesListar   = []string{"roles:listar", "roles:leer"}
	PermissionRolesCrear    = []string{"roles:crear", "roles:escribir"}
	PermissionRolesEditar   = []string{"roles:editar", "roles:escribir"}
	PermissionRolesEliminar = []string{"roles:eliminar", "roles:escribir"}
	PermissionRolesVer      = []string{"roles:ver", "roles:leer"}
	PermissionRolesPermisos = []string{"roles:permisos", "roles:leer"}
)

// Permisos para Usuarios
var (
	PermissionUsuariosListar = []string{"usuarios:listar", "usuarios:leer"}
	PermissionUsuariosCrear  = []string{"usuarios:crear", "usuarios:escribir"}
	PermissionUsuariosEditar = []string{"usuarios:editar", "usuarios:escribir"}
)

// Permisos para login
var (
	PermissionLoginLogout = []string{"login:logout", "login:cerrar_sesion"}
)

// Permisos administrativos
var (
	PermissionAdmin = []string{"admin", "administrador"}
)
