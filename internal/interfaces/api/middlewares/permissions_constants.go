package middlewares

// Permisos para Estados
var (
	PermissionEstadosListar   = []string{"estado:listar", "estado:leer"}
	PermissionEstadosCrear    = []string{"estado:escribir", "estado:crear"}
	PermissionEstadosEditar   = []string{"estado:editar", "estado:actualizar"}
	PermissionEstadosEliminar = []string{"estado:eliminar", "estado:escribir"}
	PermissionEstadosVer      = []string{"estado:ver", "estado:leer"}
)

// Permisos para Roles
var (
	PermissionRolesVer      = []string{"rol:ver", "rol:leer"}
	PermissionRolesListar   = []string{"rol:listar", "rol:leer"}
	PermissionRolesCrear    = []string{"rol:escribir", "rol:crear"}
	PermissionRolesEliminar = []string{"rol:eliminar", "rol:elimina"}
	PermissionRolesEditar   = []string{"rol:editar", "rol:actualizar"}
	PermissionRolesPermisos = []string{"rol:permisos", "rol:leer_permisos"}
)

// Permisos para Usuarios
var (
	PermissionUsuariosListar = []string{"usuarios:listar", "usuarios:leer"}
	PermissionUsuariosCrear  = []string{"usuarios:crear", "usuarios:escribir"}
	PermissionUsuariosEditar = []string{"usuarios:editar", "usuarios:escribir"}
)

// Permisos para login
var (
	PermissionLoginLogout = []string{"login:logout", "auth:logout"}
)

// Permisos administrativos
var (
	PermissionAdmin = []string{"admin", "administrador"}
)
