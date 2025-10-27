package utils

import (
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
	saltRounds int
}

func NewPasswordService(saltRounds string) *PasswordService {
	// Convertir saltRounds a int, usar 10 por defecto si no se especifica
	salt, err := strconv.Atoi(saltRounds)
	if err != nil {
		salt = 10 // Valor por defecto
	}

	return &PasswordService{
		saltRounds: salt,
	}
}

// EncodePassword encripta un password usando bcrypt
func (ps *PasswordService) EncodePassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), ps.saltRounds)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePasswords compara un password plano con uno encriptado
func (ps *PasswordService) ComparePasswords(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

// Métodos estáticos equivalentes (para mantener compatibilidad con el código TypeScript)
func EncodePassword(password string, saltRounds string) (string, error) {
	service := NewPasswordService(saltRounds)
	return service.EncodePassword(password)
}

func ComparePasswords(plainPassword, hashedPassword string) bool {
	service := NewPasswordService("10") // Salt rounds por defecto
	return service.ComparePasswords(plainPassword, hashedPassword)
}
