package auth

import (
	"encoding/json"
	"net/http"

	"api_go/internal/core/auth"
	"api_go/internal/interfaces/api/common"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := common.NewErrorResponse(400, "Solicitud inválida")
		common.WriteJSONResponse(w, response, 400)
		return
	}

	ctx := r.Context()
	response, err := h.authService.LoginUser(ctx, req)
	if err != nil {
		errorResponse := common.NewErrorResponse(401, err.Error())
		common.WriteJSONResponse(w, errorResponse, 401)
		return
	}

	common.WriteJSONResponse(w, *response, 200)
}
