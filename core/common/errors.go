package common

import "fmt"

// Custom errors
type AppError struct {
	Code    string
	Message string
	Details interface{}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewValidationError(field, message string) *AppError {
	return &AppError{
		Code:    "VALIDATION_ERROR",
		Message: fmt.Sprintf("Campo %s: %s", field, message),
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: message,
	}
}

func NewServiceError(message string) *AppError {
	return &AppError{
		Code:    "SERVICE_ERROR",
		Message: message,
	}
}
