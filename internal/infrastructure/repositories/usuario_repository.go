package repositories

import (
	"context"
	"fmt"
	"strconv"

	"api_go/internal/core/auth"
	"api_go/internal/infrastructure/database"

	"gorm.io/gorm"
)

type UsuarioRepository struct {
	dbManager *database.DBManager
}

func NewUsuarioRepository(dbManager *database.DBManager) *UsuarioRepository {
	return &UsuarioRepository{
		dbManager: dbManager,
	}
}

// FindByUsername simula: prisma.usuario.findUnique({ where: { username } })
func (r *UsuarioRepository) FindByUsername(ctx context.Context, email string) (*auth.User, error) {

	// Estructura que coincide con la BD real
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

	err := r.dbManager.FindUnique(ctx, "usuario", &usuarioDB, map[string]interface{}{
		"email": email,
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		fmt.Printf("❌ Error buscando usuario: %v\n", err)
		return nil, fmt.Errorf("error buscando usuario: %v", err)
	}

	fmt.Printf("✅ Usuario encontrado: ID=%s, NomUser=%s\n", usuarioDB.Documento, usuarioDB.NomUser)
	documentoStr := usuarioDB.Documento
	documentoInt, err := strconv.Atoi(documentoStr)
	if err != nil {
		// Manejar el error (el string no es un número válido)
		return nil, fmt.Errorf("error al convertir documento:%v", err)
	}

	return &auth.User{
		ID:           documentoInt,
		Username:     usuarioDB.NomUser,
		Email:        usuarioDB.Email,
		PasswordHash: usuarioDB.Pass,
		IDRol:        &usuarioDB.IDRol,
	}, nil
}
