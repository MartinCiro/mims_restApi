// internal/infrastructure/repositories/usuario_repository.go
package repositories

import (
	"context"
	"fmt"

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
func (r *UsuarioRepository) FindByUsername(ctx context.Context, username string) (*auth.User, error) {
	var usuarioDB struct {
		ID        int    `gorm:"column:id"`
		Nombres   string `gorm:"column:nombres"`
		Apellidos string `gorm:"column:apellidos"`
		Username  string `gorm:"column:username"`
		Pass      string `gorm:"column:pass"`
		IDRol     int    `gorm:"column:id_rol"`
		IDEstado  int    `gorm:"column:id_estado"`
		Permisos  []int  `gorm:"column:permisos"`
	}

	err := r.dbManager.FindUnique(ctx, &usuarioDB, map[string]interface{}{
		"username": username,
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error buscando usuario: %v", err)
	}

	return &auth.User{
		ID:           usuarioDB.ID,
		Username:     usuarioDB.Username,
		PasswordHash: usuarioDB.Pass,
		IDRol:        &usuarioDB.IDRol,
	}, nil
}
