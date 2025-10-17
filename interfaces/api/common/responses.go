package common

import (
	"api_go/core/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success responde con éxito (200) - Los datos van directamente en result
func Success(ctx *gin.Context, data interface{}) {
	response := common.NewResponseBody(true, http.StatusOK, data)
	ctx.JSON(http.StatusOK, response)
}

// Created responde con creación exitosa (201) - Los datos van directamente en result
func Created(ctx *gin.Context, data interface{}) {
	response := common.NewResponseBody(true, http.StatusCreated, data)
	ctx.JSON(http.StatusCreated, response)
}

// Accepted responde con aceptado (202)
func Accepted(ctx *gin.Context, data interface{}) {
	response := common.NewResponseBody(true, http.StatusAccepted, data)
	ctx.JSON(http.StatusAccepted, response)
}

// NoContent responde sin contenido (204)
func NoContent(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

// BadRequest responde con error de validación (400)
func BadRequest(ctx *gin.Context, message string) {
	result := common.MessageResult{Message: message}
	response := common.NewResponseBody(false, http.StatusBadRequest, result)
	ctx.JSON(http.StatusBadRequest, response)
}

// Unauthorized responde no autorizado (401)
func Unauthorized(ctx *gin.Context, message string) {
	result := common.MessageResult{Message: message}
	response := common.NewResponseBody(false, http.StatusUnauthorized, result)
	ctx.JSON(http.StatusUnauthorized, response)
}

// Forbidden responde prohibido (403)
func Forbidden(ctx *gin.Context, message string) {
	result := common.MessageResult{Message: message}
	response := common.NewResponseBody(false, http.StatusForbidden, result)
	ctx.JSON(http.StatusForbidden, response)
}

// NotFound responde no encontrado (404)
func NotFound(ctx *gin.Context, message string) {
	result := common.MessageResult{Message: message}
	response := common.NewResponseBody(false, http.StatusNotFound, result)
	ctx.JSON(http.StatusNotFound, response)
}

// InternalServerError responde error interno (500)
func InternalServerError(ctx *gin.Context, message string) {
	result := common.MessageResult{Message: message}
	response := common.NewResponseBody(false, http.StatusInternalServerError, result)
	ctx.JSON(http.StatusInternalServerError, response)
}

// HandleError maneja errores de la aplicación
func HandleError(ctx *gin.Context, err error) {
	if appErr, ok := err.(*common.AppError); ok {
		switch appErr.Code {
		case "NOT_FOUND":
			NotFound(ctx, appErr.Message)
		case "VALIDATION_ERROR":
			BadRequest(ctx, appErr.Message)
		case "UNAUTHORIZED":
			Unauthorized(ctx, appErr.Message)
		case "FORBIDDEN":
			Forbidden(ctx, appErr.Message)
		default:
			InternalServerError(ctx, "Error interno del servidor")
		}
		return
	}

	// Error genérico
	InternalServerError(ctx, "Error interno del servidor")
}

// SuccessWithMessage responde éxito con mensaje personalizado
func SuccessWithMessage(ctx *gin.Context, message string) {
	result := common.MessageResult{Message: message}
	response := common.NewResponseBody(true, http.StatusOK, result)
	ctx.JSON(http.StatusOK, response)
}

// SuccessWithData responde éxito con datos en formato { data: ... } - OBSOLETO, usar Success directamente
// func SuccessWithData[T any](ctx *gin.Context, data T) {
// 	result := common.DataResult[T]{Data: data}
// 	response := common.NewResponseBody(true, http.StatusOK, result)
// 	ctx.JSON(http.StatusOK, response)
// }

// SuccessWithPagination responde éxito con datos paginados
func SuccessWithPagination[T any](ctx *gin.Context, data []T, total, page, limit int) {
	result := common.PaginatedResult[T]{
		Data:  data,
		Total: total,
		Page:  page,
		Limit: limit,
	}
	response := common.NewResponseBody(true, http.StatusOK, result)
	ctx.JSON(http.StatusOK, response)
}
