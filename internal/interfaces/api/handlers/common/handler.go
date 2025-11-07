package common

import (
	"net/http"
	"time"

	"api_go/internal/interfaces/api/common"
)

// HealthHandler maneja el endpoint raíz y health check
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := common.NewErrorResponse(405, "Método no permitido")
		common.WriteJSONResponse(w, response, 405)
		return
	}

	response := common.NewSuccessResponse(map[string]string{
		"message": "Que onda Sir David",
		"status":  "running",
		"version": "1.0.0",
	})
	common.WriteJSONResponse(w, response, 200)
}

// ReadyHandler verifica que todos los servicios estén listos
func ReadyHandler(db interface{}, redis interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Aquí podrías verificar la conexión a BD y Redis
		response := common.NewSuccessResponse(map[string]interface{}{
			"status":    "ready",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		common.WriteJSONResponse(w, response, 200)
	}
}
