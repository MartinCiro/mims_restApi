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

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := common.NewErrorResponse(400, "Solicitud inválida")
		common.WriteJSONResponse(w, response, 400)
		return
	}

	// Validaciones básicas
	if req.Username == "" || req.Email == "" || req.Password == "" {
		response := common.NewErrorResponse(400, "Username, email y password son requeridos")
		common.WriteJSONResponse(w, response, 400)
		return
	}

	if len(req.Password) < 6 {
		response := common.NewErrorResponse(400, "El password debe tener al menos 6 caracteres")
		common.WriteJSONResponse(w, response, 400)
		return
	}

	ctx := r.Context()

	// Extraer usuario del contexto (si existe)
	var currentUser *auth.User
	if user, ok := ctx.Value("user").(*auth.User); ok {
		currentUser = user
	}

	response, err := h.authService.RegisterUser(ctx, req, currentUser)
	if err != nil {
		errorResponse := common.NewErrorResponse(400, err.Error())
		common.WriteJSONResponse(w, errorResponse, 400)
		return
	}

	common.WriteJSONResponse(w, response, 201)
}
