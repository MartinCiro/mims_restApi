package models

import (
	"time"
)

type Usuario struct {
	ID        int       `gorm:"primaryKey;column:id"`
	Nombres   string    `gorm:"column:nombres"`
	Apellidos string    `gorm:"column:apellidos"`
	Username  string    `gorm:"column:username;uniqueIndex"`
	Pass      string    `gorm:"column:pass"`
	IDRol     int       `gorm:"column:id_rol"`
	IDEstado  int       `gorm:"column:id_estado"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	// Relaciones
	Estado Estado `gorm:"foreignKey:IDEstado"`
	Rol    Rol    `gorm:"foreignKey:IDRol"`
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
