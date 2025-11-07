package auth

import (
	"api_go/pkg/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ExpiresAt alias para time.Time para mayor claridad
type ExpiresAt = time.Time

type Usuario struct {
	ID                *int   `json:"id,omitempty"`
	Username          string `json:"username,omitempty"`
	Email             string `json:"email"`
	IDRol             *int   `json:"id_rol,omitempty"`
	IDEstado          *int   `json:"id_estado,omitempty"`
	encryptedPassword string
}

// NewUsuario crea un nuevo usuario con password encriptado
func NewUsuario(username, password string, idRol, idEstado *int) (*Usuario, error) {
	encryptedPassword, err := encodePassword(password)
	if err != nil {
		return nil, err
	}

	return &Usuario{
		Username:          username,
		IDRol:             idRol,
		IDEstado:          idEstado,
		encryptedPassword: encryptedPassword,
	}, nil
}

// NewUsuarioFromEncrypted crea un usuario con password ya encriptado
func NewUsuarioFromEncrypted(username, encryptedPassword string, idRol, idEstado *int) *Usuario {
	return &Usuario{
		Username:          username,
		IDRol:             idRol,
		IDEstado:          idEstado,
		encryptedPassword: encryptedPassword,
	}
}

// GetEncryptedPassword retorna el password encriptado
func (u *Usuario) GetEncryptedPassword() string {
	return u.encryptedPassword
}

// ComparePassword verifica si el password plano coincide con el encriptado
func (u *Usuario) ComparePassword(plainPassword string, passwordService *utils.PasswordService) bool {
	return passwordService.ComparePasswords(plainPassword, u.encryptedPassword)
}

// encodePassword encripta un password usando bcrypt
func encodePassword(plainPassword string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
