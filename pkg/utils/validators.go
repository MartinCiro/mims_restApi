package utils

import (
	"fmt"
	"strings"

	"api_go/internal/interfaces/api/common"
)

// ValidarBlank valida que un valor no esté vacío
func ValidarBlank(valor interface{}, nombre string) error {
	if valor == nil {
		return fmt.Errorf("no se ha proporcionado %s", nombre)
	}

	switch v := valor.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("no se ha proporcionado %s", nombre)
		}
	case int, int32, int64:
		if v == 0 {
			return fmt.Errorf("no se ha proporcionado %s", nombre)
		}
	case []interface{}:
		if len(v) == 0 {
			return fmt.Errorf("no se ha proporcionado %s", nombre)
		}
		// Puedes agregar más tipos según necesites
	}

	return nil
}

// ValidarExistente valida errores de existencia (duplicados)
type ValidacionResult struct {
	OK   bool   `json:"ok"`
	Data string `json:"data,omitempty"`
}

func ValidarExistente(code string, target ...string) ValidacionResult {
	if code == "P2002" { // Código de error Prisma para duplicados
		campo := ""
		if len(target) > 0 {
			campo = target[0]
		}
		return ValidacionResult{
			OK:   false,
			Data: fmt.Sprintf("%s ya existe en la base de datos", campo),
		}
	}

	return ValidacionResult{OK: true}
}

// ValidarNoExistente valida errores de no existencia
type ValidacionNoExistenteResult struct {
	OK        bool   `json:"ok"`
	StatusCod int    `json:"status_cod,omitempty"`
	Data      string `json:"data,omitempty"`
}

func ValidarNoExistente(valor string, nombre interface{}) ValidacionNoExistenteResult {
	if valor == "P2025" { // Código de error Prisma para no encontrado
		nombreStr := fmt.Sprintf("%v", nombre)
		return ValidacionNoExistenteResult{
			OK:        false,
			StatusCod: 400,
			Data:      fmt.Sprintf("%s no existe en la base de datos", nombreStr),
		}
	}

	return ValidacionNoExistenteResult{OK: true}
}

// HandleException maneja excepciones y las convierte en respuestas HTTP
func HandleException(err error) (common.ResponseBody[interface{}], int) {
	// Error genérico
	fmt.Printf("Error interno: %v\n", err)
	return common.NewErrorResponse(500, "Error interno del servidor"), 500
}

// Capitalize convierte la primera letra a mayúscula y el resto a minúscula
func Capitalize(text string) string {
	if text == "" {
		return text
	}
	return strings.ToUpper(text[:1]) + strings.ToLower(text[1:])
}
