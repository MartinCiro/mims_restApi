package repositories

import (
	"context"
	"fmt"
	"strings"

	"api_go/internal/core/auth"
	"api_go/pkg/utils"
)

type AuthAdapter struct {
	usuarioRepo *UsuarioRepository
	rolRepo     *RolRepository
	permisoRepo *PermisoRepository
}

func NewAuthAdapter(
	usuarioRepo *UsuarioRepository,
	rolRepo *RolRepository,
	permisoRepo *PermisoRepository,
) *AuthAdapter {
	return &AuthAdapter{
		usuarioRepo: usuarioRepo,
		rolRepo:     rolRepo,
		permisoRepo: permisoRepo,
	}
}

// RetrieveUser implementa el puerto AuthPort
func (a *AuthAdapter) RetrieveUser(ctx context.Context, authData auth.AuthData) (*auth.User, error) {
	fmt.Printf("🔍 Coordinando búsqueda de usuario: %s\n", authData.Username)

	// 1. Buscar usuario
	usuario, err := a.usuarioRepo.FindByUsername(ctx, authData.Username)
	if err != nil {
		return nil, fmt.Errorf("error recuperando usuario: %v", err)
	}

	if usuario == nil {
		fmt.Printf("❌ Usuario no encontrado: %s\n", authData.Username)
		return nil, nil
	}

	fmt.Printf("✅ Usuario encontrado: ID=%d\n", usuario.ID)

	// 2. Buscar rol (opcional - si se necesita información del rol)
	if usuario.IDRol != nil {
		rol, err := a.rolRepo.FindByID(ctx, *usuario.IDRol)
		if err != nil {
			fmt.Printf("⚠️  Error obteniendo rol (continuando): %v\n", err)
			// Continuar sin información del rol
		} else if rol != nil {
			// usuario.RolNombre = rol.Nombre // Si necesitas esta info
		}
	}

	// 3. Buscar permisos (opcional)
	if usuario.IDRol != nil {
		permisos, err := a.permisoRepo.FindByRolID(ctx, *usuario.IDRol)
		if err != nil {
			fmt.Printf("⚠️  Error obteniendo permisos (continuando sin permisos): %v\n", err)
			usuario.Permisos = []string{} // Permisos vacíos
		} else {
			usuario.Permisos = permisos
			fmt.Printf("✅ Permisos encontrados: %d permisos\n", len(permisos))
		}
	}

	fmt.Printf("✅ Usuario completo recuperado: %s\n", usuario.Username)
	return usuario, nil
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
