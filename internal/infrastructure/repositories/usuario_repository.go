package repositories

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"api_go/internal/core/auth"
	"api_go/internal/infrastructure/database"

	"gorm.io/gorm"
)

// Asegurar que UsuarioRepository implemente el port
var _ auth.UserRepositoryPort = (*UsuarioRepository)(nil)

type UsuarioRepository struct {
	dbManager *database.DBManager
}

func NewUsuarioRepository(dbManager *database.DBManager) *UsuarioRepository {
	return &UsuarioRepository{
		dbManager: dbManager,
	}
}

func (r *UsuarioRepository) FindByEmailOrUsername(ctx context.Context, field, value string) (*auth.User, error) {
	fmt.Printf("🔍 Buscando usuario por %s: %s\n", field, value)

	var usuarioDB struct {
		Documento string `gorm:"column:documento"`
		Nombres   string `gorm:"column:nombres"`
		Apellido  string `gorm:"column:apellido"`
		Email     string `gorm:"column:email"`
		NomUser   string `gorm:"column:nom_user"`
		Pass      string `gorm:"column:pass"`
		IDRol     int    `gorm:"column:id_rol"`
		EstadoID  int    `gorm:"column:estado_id"`
	}

	query := r.dbManager.GetDB().WithContext(ctx).Table("usuarios")

	if field == "email" {
		query = query.Where("email = ?", value)
	} else if field == "username" {
		query = query.Where("nom_user = ?", value)
	} else {
		return nil, fmt.Errorf("campo de búsqueda inválido: %s", field)
	}

	err := query.First(&usuarioDB).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("❌ Usuario no encontrado: %s=%s\n", field, value)
			return nil, nil
		}
		fmt.Printf("❌ Error buscando usuario: %v\n", err)
		return nil, fmt.Errorf("error buscando usuario: %v", err)
	}

	fmt.Printf("✅ Usuario encontrado: Documento=%s, Email=%s, NomUser=%s\n",
		usuarioDB.Documento, usuarioDB.Email, usuarioDB.NomUser)

	documentoInt, err := strconv.Atoi(usuarioDB.Documento)
	if err != nil {
		return nil, fmt.Errorf("error al convertir documento: %v", err)
	}

	return &auth.User{
		ID:           documentoInt,
		Username:     usuarioDB.NomUser,
		Email:        usuarioDB.Email,
		PasswordHash: usuarioDB.Pass,
		IDRol:        &usuarioDB.IDRol,
		IDEstado:     &usuarioDB.EstadoID,
	}, nil
}

// FindByUsername busca usuario por email (que es el username en tu caso)
func (r *UsuarioRepository) FindByUsername(ctx context.Context, email string) (*auth.User, error) {
	fmt.Printf("🔍 Buscando usuario por email: %s\n", email)

	var usuarioDB struct {
		Documento string `gorm:"column:documento"`
		Nombres   string `gorm:"column:nombres"`
		Apellido  string `gorm:"column:apellido"`
		Email     string `gorm:"column:email"`
		NomUser   string `gorm:"column:nom_user"`
		Pass      string `gorm:"column:pass"`
		IDRol     int    `gorm:"column:id_rol"`
		EstadoID  int    `gorm:"column:estado_id"`
	}

	// Buscar por email en la tabla 'usuarios'
	err := r.dbManager.GetDB().WithContext(ctx).
		Table("usuarios").
		Where("email = ?", email).
		First(&usuarioDB).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("❌ Usuario no encontrado: %s\n", email)
			return nil, nil
		}
		fmt.Printf("❌ Error buscando usuario: %v\n", err)
		return nil, fmt.Errorf("error buscando usuario: %v", err)
	}

	fmt.Printf("✅ Usuario encontrado: Documento=%s, NomUser=%s\n", usuarioDB.Documento, usuarioDB.NomUser)

	documentoInt, err := strconv.Atoi(usuarioDB.Documento)
	if err != nil {
		return nil, fmt.Errorf("error al convertir documento: %v", err)
	}

	return &auth.User{
		ID:           documentoInt,
		Username:     usuarioDB.NomUser,
		Email:        usuarioDB.Email,
		PasswordHash: usuarioDB.Pass,
		IDRol:        &usuarioDB.IDRol,
		IDEstado:     &usuarioDB.EstadoID,
	}, nil
}

// FindByID busca un usuario por documento (que es el ID)
func (r *UsuarioRepository) FindByID(ctx context.Context, userID int) (*auth.User, error) {
	fmt.Printf("🔍 Buscando usuario por ID: %d\n", userID)

	var usuarioDB struct {
		Documento string `gorm:"column:documento"`
		Nombres   string `gorm:"column:nombres"`
		Apellido  string `gorm:"column:apellido"`
		Email     string `gorm:"column:email"`
		NomUser   string `gorm:"column:nom_user"`
		Pass      string `gorm:"column:pass"`
		IDRol     int    `gorm:"column:id_rol"`
		EstadoID  int    `gorm:"column:estado_id"`
	}

	// Buscar por documento en la tabla 'usuarios'
	err := r.dbManager.GetDB().WithContext(ctx).
		Table("usuarios").
		Where("documento = ?", strconv.Itoa(userID)).
		First(&usuarioDB).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		fmt.Printf("❌ Error buscando usuario por ID: %v\n", err)
		return nil, fmt.Errorf("error buscando usuario por ID: %v", err)
	}

	fmt.Printf("✅ Usuario encontrado por ID: Documento=%s, NomUser=%s\n", usuarioDB.Documento, usuarioDB.NomUser)

	documentoInt, err := strconv.Atoi(usuarioDB.Documento)
	if err != nil {
		return nil, fmt.Errorf("error al convertir documento a número: %v", err)
	}

	return &auth.User{
		ID:           documentoInt,
		Username:     usuarioDB.NomUser,
		Email:        usuarioDB.Email,
		PasswordHash: usuarioDB.Pass,
		IDRol:        &usuarioDB.IDRol,
		IDEstado:     &usuarioDB.EstadoID,
	}, nil
}

// CreateUser crea un nuevo usuario
func (r *UsuarioRepository) CreateUser(ctx context.Context, usuario *auth.Usuario, email string) (int, error) {
	// Generar documento único (usando timestamp)
	documento := fmt.Sprintf("%d", time.Now().UnixNano()/1000000) // Milisegundos

	dbUser := &struct {
		Documento string `gorm:"column:documento"`
		Nombres   string `gorm:"column:nombres"`
		Apellido  string `gorm:"column:apellido"`
		Email     string `gorm:"column:email"`
		NomUser   string `gorm:"column:nom_user"`
		Pass      string `gorm:"column:pass"`
		IDRol     int    `gorm:"column:id_rol"`
		EstadoID  int    `gorm:"column:estado_id"`
	}{
		Documento: documento,
		NomUser:   usuario.Username,
		Pass:      usuario.GetEncryptedPassword(),
		Email:     email,
		IDRol:     *usuario.IDRol,
		EstadoID:  *usuario.IDEstado,
		Nombres:   usuario.Username,
		Apellido:  "Usuario", // Valor por defecto
	}

	result := r.dbManager.GetDB().WithContext(ctx).
		Table("usuarios").
		Create(dbUser)

	if result.Error != nil {
		return 0, result.Error
	}

	// Convertir documento a int para el ID
	userID, err := strconv.Atoi(documento)
	if err != nil {
		return 0, fmt.Errorf("error convirtiendo documento a ID: %v", err)
	}

	fmt.Printf("✅ Usuario creado: Documento=%s, ID=%d\n", documento, userID)
	return userID, nil
}
