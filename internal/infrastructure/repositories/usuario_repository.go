package repositories

import (
	"api_go/infrastructure/database/models"
	"api_go/internal/core/auth"
	"api_go/internal/infrastructure/database"
	"api_go/pkg/logger"
	"context"
	"fmt"
	"time"

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
	logger.Info("🔍 Buscando usuario", "campo", field, "valor", value)

	var usuario models.Usuario

	query := r.dbManager.GetDB().WithContext(ctx).Model(&models.Usuario{})

	if field == "email" {
		query = query.Where("email = ?", value)
	} else if field == "username" {
		query = query.Where("nom_user = ?", value)
	} else {
		return nil, fmt.Errorf("campo de búsqueda inválido: %s", field)
	}

	err := query.First(&usuario).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Info("❌ Usuario no encontrado", "campo", field, "valor", value)
			return nil, nil
		}
		logger.Error("❌ Error buscando usuario", "error", err, "campo", field, "valor", value)
		return nil, fmt.Errorf("error buscando usuario: %v", err)
	}

	logger.Info("✅ Usuario encontrado",
		"documento", usuario.Documento,
		"email", usuario.Email,
		"nom_user", usuario.NomUser)

	return &auth.User{
		ID:           usuario.Documento,
		Username:     usuario.NomUser,
		Email:        usuario.Email,
		PasswordHash: usuario.Pass,
		IDRol:        &usuario.IDRol,
		IDEstado:     &usuario.EstadoID,
	}, nil
}

// FindByUsername busca usuario por email
func (r *UsuarioRepository) FindByUsername(ctx context.Context, email string) (*auth.User, error) {
	logger.Info("🔍 Buscando usuario por email", "email", email)

	var usuario models.Usuario

	err := r.dbManager.GetDB().WithContext(ctx).
		Model(&models.Usuario{}).
		Where("email = ?", email).
		First(&usuario).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Info("❌ Usuario no encontrado", "email", email)
			return nil, nil
		}
		logger.Error("❌ Error buscando usuario por email", "error", err, "email", email)
		return nil, fmt.Errorf("error buscando usuario: %v", err)
	}

	logger.Info("✅ Usuario encontrado",
		"documento", usuario.Documento,
		"nom_user", usuario.NomUser)

	return &auth.User{
		ID:           usuario.Documento,
		Username:     usuario.NomUser,
		Email:        usuario.Email,
		PasswordHash: usuario.Pass,
		IDRol:        &usuario.IDRol,
		IDEstado:     &usuario.EstadoID,
	}, nil
}

// FindByID busca un usuario por documento (que es el ID)
func (r *UsuarioRepository) FindByID(ctx context.Context, userID int) (*auth.User, error) {
	logger.Info("🔍 Buscando usuario por ID", "user_id", userID)

	var usuario models.Usuario

	err := r.dbManager.GetDB().WithContext(ctx).
		Model(&models.Usuario{}).
		Where("documento = ?", userID).
		First(&usuario).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Info("❌ Usuario no encontrado", "user_id", userID)
			return nil, fmt.Errorf("usuario no encontrado")
		}
		logger.Error("❌ Error buscando usuario por ID", "error", err, "user_id", userID)
		return nil, fmt.Errorf("error buscando usuario por ID: %v", err)
	}

	logger.Info("✅ Usuario encontrado por ID",
		"documento", usuario.Documento,
		"nom_user", usuario.NomUser)

	return &auth.User{
		ID:           usuario.Documento,
		Username:     usuario.NomUser,
		Email:        usuario.Email,
		PasswordHash: usuario.Pass,
		IDRol:        &usuario.IDRol,
		IDEstado:     &usuario.EstadoID,
	}, nil
}

// CreateUser crea un nuevo usuario
func (r *UsuarioRepository) CreateUser(ctx context.Context, usuario *auth.Usuario, email string) (int, error) {
	// Generar documento único (usando timestamp)
	documento := int(time.Now().UnixNano() / 1000000) // Milisegundos como int

	dbUser := &models.Usuario{
		Documento: documento,
		Nombres:   usuario.Username,
		Apellido:  "Usuario", // Valor por defecto
		Email:     email,
		NomUser:   usuario.Username,
		Pass:      usuario.GetEncryptedPassword(),
		IDRol:     *usuario.IDRol,
		EstadoID:  *usuario.IDEstado,
	}

	result := r.dbManager.GetDB().WithContext(ctx).
		Create(dbUser)

	if result.Error != nil {
		logger.Error("❌ Error creando usuario", "error", result.Error, "email", email)
		return 0, result.Error
	}

	logger.Info("✅ Usuario creado",
		"documento", documento,
		"email", email,
		"nom_user", usuario.Username)

	return documento, nil
}

// UpdateUser actualiza un usuario existente
func (r *UsuarioRepository) UpdateUser(ctx context.Context, userID int, updates map[string]interface{}) error {
	logger.Info("🔍 Actualizando usuario", "user_id", userID)

	result := r.dbManager.GetDB().WithContext(ctx).
		Model(&models.Usuario{}).
		Where("documento = ?", userID).
		Updates(updates)

	if result.Error != nil {
		logger.Error("❌ Error actualizando usuario", "error", result.Error, "user_id", userID)
		return fmt.Errorf("error actualizando usuario: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		logger.Info("❌ Usuario no encontrado para actualizar", "user_id", userID)
		return fmt.Errorf("usuario no encontrado")
	}

	logger.Info("✅ Usuario actualizado", "user_id", userID)
	return nil
}

// DeleteUser elimina un usuario por ID
func (r *UsuarioRepository) DeleteUser(ctx context.Context, userID int) error {
	logger.Info("🔍 Eliminando usuario", "user_id", userID)

	result := r.dbManager.GetDB().WithContext(ctx).
		Where("documento = ?", userID).
		Delete(&models.Usuario{})

	if result.Error != nil {
		logger.Error("❌ Error eliminando usuario", "error", result.Error, "user_id", userID)
		return fmt.Errorf("error eliminando usuario: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		logger.Info("❌ Usuario no encontrado para eliminar", "user_id", userID)
		return fmt.Errorf("usuario no encontrado")
	}

	logger.Info("✅ Usuario eliminado", "user_id", userID)
	return nil
}

// FindAll busca todos los usuarios
func (r *UsuarioRepository) FindAll(ctx context.Context) ([]models.Usuario, error) {
	logger.Info("🔍 Buscando todos los usuarios")

	var usuarios []models.Usuario

	err := r.dbManager.GetDB().WithContext(ctx).
		Find(&usuarios).Error

	if err != nil {
		logger.Error("❌ Error buscando todos los usuarios", "error", err)
		return nil, fmt.Errorf("error buscando todos los usuarios: %v", err)
	}

	logger.Info("✅ Usuarios encontrados", "cantidad", len(usuarios))
	return usuarios, nil
}

// FindByRol busca usuarios por rol
func (r *UsuarioRepository) FindByRol(ctx context.Context, rolID int) ([]models.Usuario, error) {
	logger.Info("🔍 Buscando usuarios por rol", "rol_id", rolID)

	var usuarios []models.Usuario

	err := r.dbManager.GetDB().WithContext(ctx).
		Where("id_rol = ?", rolID).
		Find(&usuarios).Error

	if err != nil {
		logger.Error("❌ Error buscando usuarios por rol", "error", err, "rol_id", rolID)
		return nil, fmt.Errorf("error buscando usuarios por rol: %v", err)
	}

	logger.Info("✅ Usuarios encontrados por rol", "rol_id", rolID, "cantidad", len(usuarios))
	return usuarios, nil
}
