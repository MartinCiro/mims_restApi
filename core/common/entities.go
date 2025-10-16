package common

// ResponseBody estructura que coincide con tu TypeScript
type ResponseBody[T any] struct {
	OK         bool `json:"ok"`
	StatusCode int  `json:"statusCode"`
	Result     T    `json:"result"`
}

// Para respuestas sin datos (solo status/message)
type MessageResult struct {
	Message string `json:"message,omitempty"`
}

// Para respuestas con datos genéricos
type DataResult[T any] struct {
	Data T `json:"data"`
}

// Para respuestas paginadas
type PaginatedResult[T any] struct {
	Data  []T `json:"data"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// Constructor para ResponseBody
func NewResponseBody[T any](ok bool, statusCode int, result T) *ResponseBody[T] {
	return &ResponseBody[T]{
		OK:         ok,
		StatusCode: statusCode,
		Result:     result,
	}
}
