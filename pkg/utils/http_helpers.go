package utils

import (
	"net/http"

	"api_go/internal/interfaces/api/common"
)

// WriteValidationError escribe un error de validación como respuesta HTTP
func WriteValidationError(w http.ResponseWriter, err error, statusCode int) {
	response, code := HandleException(err)
	common.WriteJSONResponse(w, response, code)
}

// ValidateRequired valida campos requeridos
func ValidateRequired(fields map[string]interface{}) error {
	for nombre, valor := range fields {
		if err := ValidarBlank(valor, nombre); err != nil {
			return err
		}
	}
	return nil
}
