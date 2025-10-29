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

// FindByID busca un usuario por su ID
// FindByID busca un usuario por su ID (asumiendo que userID es el documento)
func (r *UsuarioRepository) FindByID(ctx context.Context, userID int) (*auth.User, error) {
	// Estructura que coincide con la BD real
	var usuarioDB struct {
		Documento string `gorm:"column:documento"`
		Nombres   string `gorm:"column:nombres"`
		Apellido  string `gorm:"column:apellido"`
		Email     string `gorm:"column:email"`
		NomUser   string `gorm:"column:nom_user"`
		Pass      string `gorm:"column:pass"`
		IDRol     int    `gorm:"column:id_rol"`
		EstadoID  int    `gorm:"column:estado_id"` // Este campo existe en tu BD
	}

	// Buscar por documento (que es el ID del usuario)
	err := r.dbManager.FindUnique(ctx, "usuario", &usuarioDB, map[string]interface{}{
		"documento": strconv.Itoa(userID), // Convertir int a string para la búsqueda
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		fmt.Printf("❌ Error buscando usuario por ID: %v\n", err)
		return nil, fmt.Errorf("error buscando usuario por ID: %v", err)
	}

	fmt.Printf("✅ Usuario encontrado por ID: Documento=%s, NomUser=%s\n", usuarioDB.Documento, usuarioDB.NomUser)

	// Convertir documento string a int para el ID
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
		IDEstado:     &usuarioDB.EstadoID, // Agregar este campo que falta
	}, nil
}
