package repositories

import (
	"context"
	"fmt"

	"api_go/internal/core/auth"
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

	fmt.Printf("✅ Usuario base encontrado: ID=%d\n", usuario.ID)

	// 2. Buscar permisos (si tiene rol)
	if usuario.IDRol != nil {
		permisos, err := a.permisoRepo.FindByRolID(ctx, *usuario.IDRol)
		if err != nil {
			fmt.Printf("⚠️  Error obteniendo permisos (continuando sin permisos): %v\n", err)
			usuario.Permisos = []string{}
		} else {
			usuario.Permisos = permisos
			fmt.Printf("✅ Permisos asignados: %d permisos\n", len(permisos))
		}
	}

	// 3. Buscar información del rol (opcional)
	if usuario.IDRol != nil {
		rol, err := a.rolRepo.FindByID(ctx, *usuario.IDRol)
		if err != nil {
			fmt.Printf("⚠️  Error obteniendo rol (continuando): %v\n", err)
		} else if rol != nil {
			// usuario.RolNombre = rol.Nombre // Si necesitas esta info
			fmt.Printf("✅ Información de rol obtenida: %s\n", rol.Nombre)
		}
	}

	fmt.Printf("✅ Usuario completo recuperado: %s\n", usuario.Username)
	return usuario, nil
}
