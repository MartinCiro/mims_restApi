package common

import (
	"encoding/json"
	"net/http"
)

// ResponseBody representa la estructura estándar de respuesta
type ResponseBody[T any] struct {
	Success bool   `json:"ok"`
	Code    int    `json:"statusCode"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"result,omitempty"`
}

// NewResponseBody crea una nueva respuesta
func NewResponseBody[T any](success bool, code int, data T) ResponseBody[T] {
	return ResponseBody[T]{
		Success: success,
		Code:    code,
		Data:    data,
	}
}

// NewSuccessResponse crea una respuesta exitosa
func NewSuccessResponse[T any](data T) ResponseBody[T] {
	return ResponseBody[T]{
		Success: true,
		Code:    200,
		Message: "Operación exitosa",
		Data:    data,
	}
}

// NewErrorResponse crea una respuesta de error
func NewErrorResponse(code int, message string) ResponseBody[interface{}] {
	return ResponseBody[interface{}]{
		Success: false,
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

// WriteJSONResponse escribe una respuesta JSON estandarizada
func WriteJSONResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
