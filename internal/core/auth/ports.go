package auth

import "context"

type AuthPort interface {
	RetrieveUser(ctx context.Context, authData AuthData) (*User, error)
	RetrieveUserByID(ctx context.Context, userID int) (*User, error)
}

type UserRepositoryPort interface {
	CreateUser(ctx context.Context, usuario *Usuario, email string) (int, error)
}

type RolRepositoryPort interface {
	FindRolIDByName(ctx context.Context, nombre string) (int, error)
}

type EstadoRepositoryPort interface {
	FindEstadoIDByName(ctx context.Context, nombre string) (int, error)
}

type AuthData struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type User struct {
	ID           int      `json:"id"`
	Username     string   `json:"username,omitempty"`
	Email        string   `json:"email"`
	PasswordHash string   `json:"password_hash,omitempty"`
	IDRol        *int     `json:"id_rol,omitempty"`
	IDEstado     *int     `json:"id_estado,omitempty"`
	Permisos     []string `json:"permisos,omitempty"`
	RolNombre    string   `json:"rol_nombre,omitempty"`
}
