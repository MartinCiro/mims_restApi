package models

import (
	"time"
)

type Usuario struct {
	// ✅ CORREGIR: Mapear correctamente a las columnas reales
	Documento string `gorm:"primaryKey;column:documento;type:text"`
	Nombres   string `gorm:"column:nombres;type:text"`
	Apellidos string `gorm:"column:apellido;type:text"`
	Email     string `gorm:"column:email;type:text"`
	Username  string `gorm:"column:nom_user;type:text"`
	Pass      string `gorm:"column:pass;type:text"`
	IDRol     int64  `gorm:"column:id_rol;type:bigint"`
	IDEstado  int64  `gorm:"column:estado_id;type:bigint"`

	// Campos adicionales que existen en tu BD
	InfoPerfil      string    `gorm:"column:info_perfil;type:text"`
	NumContacto     string    `gorm:"column:num_contacto;type:text"`
	FechaNacimiento time.Time `gorm:"column:fecha_nacimiento"`
	FechaRegistro   time.Time `gorm:"column:fecha_registro"`

	// Relaciones
	Estado Estado `gorm:"foreignKey:IDEstado;references:ID"`
	Rol    Rol    `gorm:"foreignKey:IDRol;references:ID"`
}

type Estado struct {
	ID     int    `gorm:"primaryKey;column:id"`
	Nombre string `gorm:"column:nombre;uniqueIndex"`
}

type Rol struct {
	ID     int    `gorm:"primaryKey;column:id"`
	Nombre string `gorm:"column:nombre;uniqueIndex"`
}

type Permiso struct {
	ID     int    `gorm:"primaryKey;column:id"`
	Nombre string `gorm:"column:nombre;uniqueIndex"`
}

type RolXPermiso struct {
	IDRol     int `gorm:"primaryKey;column:id_rol"`
	IDPermiso int `gorm:"primaryKey;column:id_permiso"`
}
