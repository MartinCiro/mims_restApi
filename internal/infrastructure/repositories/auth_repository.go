package repositories

import (
	"context"
	"fmt"
	"strings"

	"api_go/internal/core/auth"
	"api_go/pkg/utils"

	"gorm.io/gorm"
)

type AuthAdapter struct {
	db *gorm.DB
}

func NewAuthAdapter(db *gorm.DB) *AuthAdapter {
	return &AuthAdapter{
		db: db,
	}
}

// RetrieveUser implementa el puerto AuthPort
func (a *AuthAdapter) RetrieveUser(ctx context.Context, authData auth.AuthData) (*auth.User, error) {
	type UsuarioDB struct {
		ID        int    `gorm:"column:id"`
		Nombres   string `gorm:"column:nombres"`
		Apellidos string `gorm:"column:apellidos"`
		Username  string `gorm:"column:username"`
		Pass      string `gorm:"column:pass"`
		IDRol     int    `gorm:"column:id_rol"`
	}

	type EstadoDB struct {
		Nombre string `gorm:"column:nombre"`
	}

	type RolDB struct {
		Nombre string `gorm:"column:nombre"`
	}

	type RolXPermisoDB struct {
		Permiso struct {
			Nombre string `gorm:"column:nombre"`
		} `gorm:"embedded;embeddedPrefix:permiso_"`
	}

	var usuarioDB UsuarioDB

	// Consultar usuario con sus relaciones
	err := a.db.WithContext(ctx).
		Table("usuario").
		Select(`
			usuario.id,
			usuario.nombres,
			usuario.apellidos, 
			usuario.username,
			usuario.pass,
			usuario.id_rol
		`).
		Joins("LEFT JOIN estado ON usuario.id_estado = estado.id").
		Joins("LEFT JOIN rol ON usuario.id_rol = rol.id").
		Where("usuario.username = ?", authData.Username).
		First(&usuarioDB).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Usuario no encontrado, no es error
		}

		// Manejar errores de base de datos
		return nil, a.handleDBError(err, authData.Username)
	}

	// Obtener estado del usuario
	var estadoDB EstadoDB
	err = a.db.WithContext(ctx).
		Table("estado").
		Select("nombre").
		Joins("INNER JOIN usuario ON estado.id = usuario.id_estado").
		Where("usuario.id = ?", usuarioDB.ID).
		First(&estadoDB).Error
	if err != nil {
		fmt.Printf("Error obteniendo estado: %v\n", err)
	}

	// Obtener rol y permisos
	var rolDB RolDB
	err = a.db.WithContext(ctx).
		Table("rol").
		Select("nombre").
		Where("id = ?", usuarioDB.IDRol).
		First(&rolDB).Error
	if err != nil {
		fmt.Printf("Error obteniendo rol: %v\n", err)
	}

	// Obtener permisos del rol
	var rolXPermisosDB []RolXPermisoDB
	err = a.db.WithContext(ctx).
		Table("rol_x_permiso").
		Select("permiso.nombre as permiso_nombre").
		Joins("INNER JOIN permiso ON rol_x_permiso.id_permiso = permiso.id").
		Where("rol_x_permiso.id_rol = ?", usuarioDB.IDRol).
		Find(&rolXPermisosDB).Error
	if err != nil {
		fmt.Printf("Error obteniendo permisos: %v\n", err)
	}

	// Extraer nombres de permisos
	permisos := make([]string, len(rolXPermisosDB))
	for i, rp := range rolXPermisosDB {
		permisos[i] = rp.Permiso.Nombre
	}

	// Mapear a entidad del core
	return &auth.User{
		ID:           usuarioDB.ID,
		Username:     usuarioDB.Username,
		PasswordHash: usuarioDB.Pass,
		IDRol:        &usuarioDB.IDRol,
		Permisos:     permisos,
	}, nil
}

// handleDBError maneja errores de base de datos
func (a *AuthAdapter) handleDBError(err error, username string) error {
	// Convertir error de GORM a códigos similares a Prisma
	errStr := err.Error()

	// Detectar errores de duplicado (similar a P2002 de Prisma)
	if strings.Contains(errStr, "duplicate") || strings.Contains(errStr, "Duplicate") {
		validacion := utils.ValidarExistente("P2002", username)
		if !validacion.OK {
			return fmt.Errorf("status_cod:409, data:%s", validacion.Data)
		}
	}

	// Error genérico
	return fmt.Errorf("status_cod:400, data:%s", "Ocurrió un error consultando el usuario")
}
