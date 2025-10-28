// internal/interfaces/api/common/responses.go
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
