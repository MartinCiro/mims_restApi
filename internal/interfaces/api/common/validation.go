// internal/interfaces/api/common/validation.go
package common

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateRequest valida un struct usando validator
func ValidateRequest(s interface{}) map[string]string {
	err := validate.Struct(s)

	if err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			switch err.Tag() {
			case "required":
				errors[field] = "Este campo es obligatorio"
			case "min":
				if err.Field() == "ID" {
					errors[field] = "El ID debe ser mayor a 0"
				} else {
					errors[field] = "El valor es demasiado corto"
				}
			case "max":
				errors[field] = "El valor excede el máximo permitido"
			default:
				errors[field] = "Valor inválido"
			}
		}
		return errors
	}
	return nil
}

// WriteValidationErrors escribe errores de validación como respuesta JSON
func WriteValidationErrors(w http.ResponseWriter, errors map[string]string, statusCode int) {
	// Crear la estructura de error que coincide con tu formato
	errorResponse := map[string]interface{}{
		"details": errors,
	}

	response := ResponseBody[map[string]interface{}]{
		Success: false,
		Code:    statusCode,
		Data:    errorResponse,
	}

	WriteJSONResponse(w, response, statusCode)
}

// Nueva función para errores simples
func WriteSimpleError(w http.ResponseWriter, message string, statusCode int) {
	response := ResponseBody[string]{
		Success: false,
		Code:    statusCode,
		Data:    message,
	}
	WriteJSONResponse(w, response, statusCode)
}
