package auth

import "context"

// AuthPort define el contrato para la autenticación
type AuthPort interface {
	RetrieveUser(ctx context.Context, authData AuthData) (*User, error)
	RetrieveUserByID(ctx context.Context, userID int) (*User, error)
}

// AuthData contiene los datos necesarios para la autenticación
type AuthData struct {
	Username string `json:"email"`
	// Puedes agregar más campos si son necesarios, como email, etc.
}

// User representa la entidad usuario devuelta por el puerto
type User struct {
	ID           int      `json:"id"`
	Username     string   `json:"username,omitempty"`
	Email        string   `json:"email"`
	PasswordHash string   `json:"password_hash,omitempty"`
	IDRol        *int     `json:"id_rol,omitempty"`
	IDEstado     *int     `json:"id_estado,omitempty"`
	Permisos     []string `json:"permisos,omitempty"`
}
