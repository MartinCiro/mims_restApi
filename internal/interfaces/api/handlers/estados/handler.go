package estados

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_go/internal/core/estados"
	"api_go/internal/interfaces/api/common"
	"api_go/pkg/utils"
)

type EstadosHandler struct {
	estadoService *estados.EstadoService
}

func NewEstadosHandler(estadoService *estados.EstadoService) *EstadosHandler {
	return &EstadosHandler{
		estadoService: estadoService,
	}
}

func (h *EstadosHandler) ObtenerEstados(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	estados, err := h.estadoService.ObtenerEstados(ctx)
	if err != nil {
		utils.WriteValidationError(w, err, 500)
		return
	}

	response := common.NewSuccessResponse(estados)
	common.WriteJSONResponse(w, response, 200)
}

func (h *EstadosHandler) ObtenerEstadoXid(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := common.NewErrorResponse(400, "ID inválido")
		common.WriteJSONResponse(w, response, 400)
		return
	}

	ctx := r.Context()
	estadoData := estados.EstadoDataXid{ID: id}
	estado, err := h.estadoService.ObtenerEstadoXid(ctx, estadoData)
	if err != nil {
		utils.WriteValidationError(w, err, 500)
		return
	}

	if estado == nil {
		response := common.NewErrorResponse(404, "Estado no encontrado")
		common.WriteJSONResponse(w, response, 404)
		return
	}

	response := common.NewSuccessResponse(estado)
	common.WriteJSONResponse(w, response, 200)
}

func (h *EstadosHandler) CrearEstado(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Nombre      string  `json:"nombre"`
		Descripcion *string `json:"descripcion,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := common.NewErrorResponse(400, "Solicitud inválida")
		common.WriteJSONResponse(w, response, 400)
		return
	}

	ctx := r.Context()
	estadoData := estados.EstadoData{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
	}

	estado, err := h.estadoService.CrearEstado(ctx, estadoData)
	if err != nil {
		utils.WriteValidationError(w, err, 400)
		return
	}

	response := common.NewSuccessResponse(estado)
	common.WriteJSONResponse(w, response, 201)
}

func (h *EstadosHandler) ActualizarEstado(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := common.NewErrorResponse(400, "ID inválido")
		common.WriteJSONResponse(w, response, 400)
		return
	}

	var req struct {
		Nombre      string  `json:"nombre"`
		Descripcion *string `json:"descripcion,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := common.NewErrorResponse(400, "Solicitud inválida")
		common.WriteJSONResponse(w, response, 400)
		return
	}

	ctx := r.Context()
	estadoData := estados.EstadoDataUpdate{
		ID:          id,
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
	}

	estado, err := h.estadoService.UpEstado(ctx, estadoData)
	if err != nil {
		utils.WriteValidationError(w, err, 400)
		return
	}

	response := common.NewSuccessResponse(estado)
	common.WriteJSONResponse(w, response, 200)
}

func (h *EstadosHandler) EliminarEstado(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := common.NewErrorResponse(400, "ID inválido")
		common.WriteJSONResponse(w, response, 400)
		return
	}

	ctx := r.Context()
	estadoData := estados.EstadoDataXid{ID: id}
	err = h.estadoService.DelEstado(ctx, estadoData)
	if err != nil {
		utils.WriteValidationError(w, err, 400)
		return
	}

	response := common.NewSuccessResponse(map[string]string{
		"message": "Estado eliminado correctamente",
	})
	common.WriteJSONResponse(w, response, 200)
}
