package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type Usuario struct {
	ID                *int   `json:"id,omitempty"`
	Username          string `json:"username"`
	Email             string `json:"email,omitempty"` // Agregado para flexibilidad
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

// comparePasswords compara un password plano con uno encriptado
func comparePasswords(plainPassword, encryptedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(plainPassword))
	return err == nil
}
