// internal/core/auth/ports.go
package auth

import "context"

// AuthPort define el contrato para la autenticación
type AuthPort interface {
	RetrieveUser(ctx context.Context, authData AuthData) (*User, error)
}

// AuthData contiene los datos necesarios para la autenticación
type AuthData struct {
	Username string `json:"username"`
	// Puedes agregar más campos si son necesarios, como email, etc.
}

// User representa la entidad usuario devuelta por el puerto
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email,omitempty"`
	PasswordHash string `json:"password_hash,omitempty"`
	IDRol        *int   `json:"id_rol,omitempty"`
	IDEstado     *int   `json:"id_estado,omitempty"`
}
