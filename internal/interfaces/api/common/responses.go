package common

import (
	"encoding/json"
	"net/http"
)

// ResponseBody es un CONCERN HTTP, vive solo en interfaces
type ResponseBody[T any] struct {
	Success bool `json:"ok"`
	Code    int  `json:"statusCode"`
	Data    T    `json:"result"`
}

func NewResponseBody[T any](success bool, code int, data T) ResponseBody[T] {
	return ResponseBody[T]{
		Success: success,
		Code:    code,
		Data:    data,
	}
}

func NewSuccessResponse[T any](data T) ResponseBody[T] {
	return ResponseBody[T]{
		Success: true,
		Code:    200,
		Data:    data,
	}
}

func NewErrorResponse(code int, message string) ResponseBody[string] {
	return ResponseBody[string]{
		Success: false,
		Code:    code,
		Data:    message,
	}
}

func WriteJSONResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Helpers para Gin (si los necesitas)
func Success(w http.ResponseWriter, data interface{}) {
	response := NewSuccessResponse(data)
	WriteJSONResponse(w, response, http.StatusOK)
}

func BadRequest(w http.ResponseWriter, message string) {
	response := NewErrorResponse(http.StatusBadRequest, message)
	WriteJSONResponse(w, response, http.StatusBadRequest)
}

func Unauthorized(w http.ResponseWriter, message string) {
	response := NewErrorResponse(http.StatusUnauthorized, message)
	WriteJSONResponse(w, response, http.StatusUnauthorized)
}

func Forbidden(w http.ResponseWriter, message string) {
	response := NewErrorResponse(http.StatusForbidden, message)
	WriteJSONResponse(w, response, http.StatusForbidden)
}

func NotFound(w http.ResponseWriter, message string) {
	response := NewErrorResponse(http.StatusNotFound, message)
	WriteJSONResponse(w, response, http.StatusNotFound)
}

func InternalServerError(w http.ResponseWriter, message string) {
	response := NewErrorResponse(http.StatusInternalServerError, message)
	WriteJSONResponse(w, response, http.StatusInternalServerError)
}
