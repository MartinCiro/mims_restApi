package usuarios

import (
	"api_go/internal/core/usuarios"
	"api_go/internal/infrastructure/database/models"
	"api_go/internal/infrastructure/redis"
	"api_go/pkg/logger"
	"api_go/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UsuariosAdapter struct {
	db           *gorm.DB
	redisService *redis.Cache
}

func NewUsuariosAdapter(db *gorm.DB, redisService *redis.Cache) *UsuariosAdapter {
	return &UsuariosAdapter{
		db:           db,
		redisService: redisService,
	}
}

func (a *UsuariosAdapter) toUsuarioEntityConRelaciones(usuarioDB *models.Usuario) *usuarios.Usuario {
	// Obtener nombres desde las relaciones cargadas con Preload
	rolNombre := a.obtenerNombreRolDesdeRelacion(usuarioDB)
	estadoNombre := a.obtenerNombreEstadoDesdeRelacion(usuarioDB)

	return &usuarios.Usuario{
		Documento:          usuarioDB.Documento,
		Nombres:            usuarioDB.Nombres,
		Apellido:           usuarioDB.Apellido,
		NombreCompleto:     usuarioDB.Nombres + " " + usuarioDB.Apellido,
		Email:              usuarioDB.Email,
		Username:           usuarioDB.NomUser,
		RolID:              usuarioDB.IDRol,
		EstadoID:           usuarioDB.EstadoID,
		RolNombre:          rolNombre,
		EstadoNombre:       estadoNombre,
		FechaRegistro:      usuarioDB.FechaRegistro,
		FechaActualizacion: time.Now(),
	}
}

func (a *UsuariosAdapter) obtenerNombreRolDesdeRelacion(usuarioDB *models.Usuario) string {
	// ✅ Si el Preload funcionó, úsalo. Si no, devuelve vacío (el Core o la BD deben garantizar la integridad).
	return usuarioDB.Rol.NombreRol
}

func (a *UsuariosAdapter) obtenerNombreEstadoDesdeRelacion(usuarioDB *models.Usuario) string {
	return usuarioDB.Estado.NombreEstado
}

// ObtenerUsuarios implementa el puerto UsuariosPort
func (a *UsuariosAdapter) ObtenerUsuarios(ctx context.Context) ([]usuarios.Usuario, error) {
	cacheKey := "usuarios:leer"

	// Intentar obtener del cache
	cachedUsuarios, err := a.redisService.Get(ctx, cacheKey)
	if err == nil && cachedUsuarios != "" {
		var usuariosCache []usuarios.Usuario
		if err := json.Unmarshal([]byte(cachedUsuarios), &usuariosCache); err == nil {
			return usuariosCache, nil
		}
	}

	// Consultar base de datos CARGANDO RELACIONES
	var usuariosDB []models.Usuario
	err = a.db.WithContext(ctx).
		Preload("Rol").    // Cargar relación Rol
		Preload("Estado"). // Cargar relación Estado
		Find(&usuariosDB).Error

	if err != nil {
		return nil, a.handleQueryError(err, "consultando usuarios")
	}

	if len(usuariosDB) == 0 {
		return []usuarios.Usuario{}, nil
	}

	// Mapear a entidades del core
	usuariosList := make([]usuarios.Usuario, len(usuariosDB))
	for i, usuarioDB := range usuariosDB {
		usuariosList[i] = *a.toUsuarioEntityConRelaciones(&usuarioDB) // Cambia a esta función
	}

	// Guardar en cache
	usuariosJSON, err := json.Marshal(usuariosList)
	if err == nil {
		a.redisService.Set(ctx, cacheKey, string(usuariosJSON), 1800)
	}

	return usuariosList, nil
}

// ObtenerUsuarioXid implementa el puerto UsuariosPort
func (a *UsuariosAdapter) ObtenerUsuarioXid(ctx context.Context, usuarioData usuarios.UsuarioDataXid) (*usuarios.Usuario, error) {
	documento := fmt.Sprintf("%s", usuarioData.Documento)
	logger.Warn("Consultando usuario con documento", "documento", usuarioData.Documento)

	cacheKey := fmt.Sprintf("usuars:%s", usuarioData.Documento)

	// Intentar obtener del cache
	cachedUsuario, err := a.redisService.Get(ctx, cacheKey)
	if err == nil && cachedUsuario != "" {
		var usuarioCache usuarios.Usuario
		if err := json.Unmarshal([]byte(cachedUsuario), &usuarioCache); err == nil {
			return &usuarioCache, nil
		}
	}

	// SOLUCIÓN: Usar una consulta más explícita
	var usuarioDB models.Usuario

	// Opción 1: Usar Find con Limit(1)
	logger.Warn("Este es doc %d", documento)
	err = a.db.WithContext(ctx).
		Preload("Rol").
		Preload("Estado").
		Where("documento = ?", documento).
		Limit(1).
		Find(&usuarioDB).Error

	if err != nil {
		return nil, a.handleQueryError(err, "consultando usuario")
	}

	// Verificar si se encontró algún registro
	if usuarioDB.Documento == "" {
		return nil, nil
	}

	usuarioEntity := a.toUsuarioEntityConRelaciones(&usuarioDB)

	// Guardar en cache
	usuarioJSON, err := json.Marshal(usuarioEntity)
	if err == nil {
		a.redisService.Set(ctx, cacheKey, string(usuarioJSON), 1800)
	}

	return usuarioEntity, nil
}

// ActualizarUsuario implementa el puerto UsuariosPort
func (a *UsuariosAdapter) ActualizarUsuario(ctx context.Context, usuarioData usuarios.UsuarioDataUpdate) (*usuarios.Usuario, error) {
	// Verificar si el usuario existe
	var usuarioExistente models.Usuario
	err := a.db.WithContext(ctx).Where("documento = ?", usuarioData.Documento).First(&usuarioExistente).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("el usuario solicitado no existe")
		}
		return nil, a.handleQueryError(err, "verificando usuario existente")
	}

	// Preparar updates
	updates := make(map[string]interface{})

	// Campos básicos
	if usuarioData.Username != nil {
		updates["NomUser"] = *usuarioData.Username
	}
	if usuarioData.Email != nil {
		updates["email"] = *usuarioData.Email
	}
	if usuarioData.Nombres != nil {
		updates["nombres"] = *usuarioData.Nombres
	}
	if usuarioData.Apellido != nil {
		updates["apellido"] = *usuarioData.Apellido
	}

	// Campos de perfil
	if usuarioData.InfoPerfil != nil {
		updates["info_perfil"] = *usuarioData.InfoPerfil
	}
	if usuarioData.NumContacto != nil {
		updates["num_contacto"] = *usuarioData.NumContacto
	}
	if usuarioData.FechaNacimiento != nil {
		updates["fecha_nacimiento"] = *usuarioData.FechaNacimiento
	}

	// Campos de relaciones
	if usuarioData.RolID != nil {
		updates["IDRol"] = *usuarioData.RolID
	}
	if usuarioData.EstadoID != nil {
		updates["EstadoID"] = *usuarioData.EstadoID // ← Cambiado a EstadoID
	}

	// Agregar fecha de actualización (si existe en tu tabla)
	// updates["fecha_actualizacion"] = time.Now()

	// Si no hay campos para actualizar, retornar el usuario existente
	if len(updates) == 0 {
		return a.toUsuarioEntityConRelaciones(&usuarioExistente), nil
	}

	// Actualizar en base de datos
	err = a.db.WithContext(ctx).Model(&models.Usuario{}).
		Where("documento = ?", usuarioData.Documento).
		Updates(updates).Error

	if err != nil {
		return nil, a.handleUpdateError(err, usuarioData.Username, usuarioData.Email)
	}

	// Obtener usuario actualizado CON LOS NOMBRES CORRECTOS DE RELACIONES
	var usuarioActualizado models.Usuario
	err = a.db.WithContext(ctx).
		Preload("Rol").
		Preload("Estado").
		Where("documento = ?", usuarioData.Documento).
		First(&usuarioActualizado).Error

	if err != nil {
		return nil, a.handleQueryError(err, "obteniendo usuario actualizado")
	}

	// Actualizar cache
	a.actualizarCacheUsuario(ctx, usuarioData.Documento, &usuarioActualizado)

	return a.toUsuarioEntityConRelaciones(&usuarioActualizado), nil
}

// EliminarUsuario implementa el puerto UsuariosPort
func (a *UsuariosAdapter) EliminarUsuario(ctx context.Context, usuarioData usuarios.UsuarioDataXid) error {
	// Actualizar cache de lista
	//a.actualizarCacheListaEliminado(ctx, strconv(usuarioData.Documento))

	// Eliminar de base de datos
	result := a.db.WithContext(ctx).Where("documento = ?", usuarioData.Documento).Delete(&models.Usuario{})
	if result.Error != nil {
		return a.handleDeleteError(result.Error, usuarioData.Documento)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("el usuario con ID %s no existe", usuarioData.Documento)
	}

	// Limpiar cache individual
	a.redisService.Delete(ctx, fmt.Sprintf("usuario:%s", usuarioData.Documento))

	return nil
}

// CambiarPassword implementa el puerto UsuariosPort
func (a *UsuariosAdapter) CambiarPassword(ctx context.Context, cambiarPasswordData usuarios.CambiarPasswordData) error {
	// Hash de la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cambiarPasswordData.NuevoPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error procesando nueva contraseña")
	}

	// Actualizar contraseña
	err = a.db.WithContext(ctx).Model(&models.Usuario{}).
		Where("documento = ?", cambiarPasswordData.ID).
		Updates(map[string]interface{}{
			"password": string(hashedPassword),
		}).Error

	if err != nil {
		return a.handleUpdateError(err, nil, nil)
	}

	// Limpiar cache del usuario
	a.redisService.Delete(ctx, fmt.Sprintf("usuario:%d", cambiarPasswordData.ID))

	return nil
}

// ObtenerUsuarioPorUsername implementa el puerto UsuariosPort
func (a *UsuariosAdapter) ObtenerUsuarioPorUsername(ctx context.Context, username string) (*usuarios.Usuario, error) {
	var usuarioDB models.Usuario
	err := a.db.WithContext(ctx).
		Where("username = ?", username).
		First(&usuarioDB).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, a.handleQueryError(err, "consultando usuario por username")
	}

	return a.toUsuarioEntityConRelaciones(&usuarioDB), nil
}

// ObtenerUsuarioPorEmail implementa el puerto UsuariosPort
func (a *UsuariosAdapter) ObtenerUsuarioPorEmail(ctx context.Context, email string) (*usuarios.Usuario, error) {
	var usuarioDB models.Usuario
	err := a.db.WithContext(ctx).
		Where("email = ?", email).
		First(&usuarioDB).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, a.handleQueryError(err, "consultando usuario por email")
	}

	return a.toUsuarioEntityConRelaciones(&usuarioDB), nil
}

// ObtenerUsuariosConRelaciones implementa el puerto UsuariosPort
func (a *UsuariosAdapter) ObtenerUsuariosConRelaciones(ctx context.Context) ([]usuarios.UsuarioConRelaciones, error) {
	var usuariosDB []models.Usuario
	err := a.db.WithContext(ctx).
		Preload("Rol").
		Preload("Estado").
		Find(&usuariosDB).Error

	if err != nil {
		return nil, a.handleQueryError(err, "consultando usuarios con relaciones")
	}

	usuariosList := make([]usuarios.UsuarioConRelaciones, len(usuariosDB))
	for i, usuarioDB := range usuariosDB {
		usuariosList[i] = *a.toUsuarioConRelacionesEntity(&usuarioDB)
	}

	return usuariosList, nil
}

// ObtenerUsuarioConRelaciones implementa el puerto UsuariosPort
func (a *UsuariosAdapter) ObtenerUsuarioConRelaciones(ctx context.Context, usuarioData usuarios.UsuarioDataXid) (*usuarios.UsuarioConRelaciones, error) {
	var usuarioDB models.Usuario
	err := a.db.WithContext(ctx).
		Preload("Rol").
		Preload("Estado").
		Where("documento = ?", usuarioData.Documento).
		First(&usuarioDB).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, a.handleQueryError(err, "consultando usuario con relaciones")
	}

	return a.toUsuarioConRelacionesEntity(&usuarioDB), nil
}

// Mapeo de DB a Entity con relaciones
func (a *UsuariosAdapter) toUsuarioConRelacionesEntity(usuarioDB *models.Usuario) *usuarios.UsuarioConRelaciones {
	rolNombre := ""

	if usuarioDB.IDRol != 0 {
		rolNombre = usuarioDB.Rol.NombreRol
	}

	estadoNombre := ""
	if usuarioDB.Estado.ID != 0 {
		estadoNombre = usuarioDB.Estado.NombreEstado
	}

	return &usuarios.UsuarioConRelaciones{
		Username:  usuarioDB.NomUser,
		Email:     usuarioDB.Email,
		RolID:     usuarioDB.IDRol,
		EstadoID:  usuarioDB.EstadoID,
		RolNombre: rolNombre,
		Estado:    estadoNombre,
	}
}

// Métodos auxiliares para cache
func (a *UsuariosAdapter) actualizarCacheUsuario(ctx context.Context, userID string, usuarioDB *models.Usuario) {
	// Actualizar cache individual
	usuarioEntity := a.toUsuarioEntityConRelaciones(usuarioDB)
	usuarioJSON, err := json.Marshal(usuarioEntity)
	if err == nil {
		a.redisService.Set(ctx, fmt.Sprintf("usuario:%s", userID), string(usuarioJSON), 1800)
	}

	// Actualizar cache de lista
	a.actualizarCacheLista(ctx, userID, usuarioEntity)
}

func (a *UsuariosAdapter) actualizarCacheLista(ctx context.Context, userID string, usuarioEntity *usuarios.Usuario) {
	cachedUsuarios, err := a.redisService.Get(ctx, "usuarios:leer")
	if err == nil && cachedUsuarios != "" {
		var usuariosCache []usuarios.Usuario
		if err := json.Unmarshal([]byte(cachedUsuarios), &usuariosCache); err == nil {
			// Actualizar usuario en la lista
			for i, u := range usuariosCache {
				if u.Documento == userID {
					usuariosCache[i] = *usuarioEntity
					break
				}
			}
			updatedJSON, _ := json.Marshal(usuariosCache)
			a.redisService.Set(ctx, "usuarios:leer", string(updatedJSON), 1800)
		}
	}
}

func (a *UsuariosAdapter) handleQueryError(err error, operation string) error {
	return fmt.Errorf("Ocurrió un error %s", operation)
}

func (a *UsuariosAdapter) handleDeleteError(err error, id string) error {
	errStr := err.Error()
	logger.Error("Error eliminando usuario", "id", id)

	if strings.Contains(errStr, "foreign") || strings.Contains(errStr, "constraint") {
		return fmt.Errorf("No se puede eliminar el usuario porque tiene registros asociados")
	}

	return fmt.Errorf("Ocurrió un error eliminando el usuario")
}

func (a *UsuariosAdapter) handleUpdateError(err error, username, email *string) error {
	errStr := err.Error()

	if strings.Contains(errStr, "duplicate") || strings.Contains(errStr, "Duplicate") {
		if username != nil && strings.Contains(errStr, "username") {
			validacion := utils.ValidarExistente("P2002", *username)
			if !validacion.OK {
				return fmt.Errorf("status_cod:409, data:%s", validacion.Data)
			}
		}
		if email != nil && strings.Contains(errStr, "email") {
			validacion := utils.ValidarExistente("P2002", *email)
			if !validacion.OK {
				return fmt.Errorf("status_cod:409, data:%s", validacion.Data)
			}
		}
	}
	return fmt.Errorf("Ocurrió un error actualizando el usuario")
}
