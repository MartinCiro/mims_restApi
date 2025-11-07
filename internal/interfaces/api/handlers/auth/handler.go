package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"api_go/internal/core/auth"
	"api_go/internal/infrastructure/cookies"
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
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := common.NewErrorResponse(400, "Solicitud inválida")
		common.WriteJSONResponse(w, response, 400)
		return
	}

	if req.Email == "" {
		response := common.NewErrorResponse(400, "El email es requerido")
		common.WriteJSONResponse(w, response, 400)
		return
	}

	ctx := r.Context()

	loginReq := auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	authResponse, signedCookie, expiresAt, err := h.authService.LoginUser(ctx, loginReq)
	if err != nil {
		errorResponse := common.NewErrorResponse(401, err.Error())
		common.WriteJSONResponse(w, errorResponse, 401)
		return
	}

	if signedCookie != "" {
		cookies.SetAuthCookie(w, signedCookie, expiresAt)
		fmt.Printf("🎯 COOKIE SET IN LOGIN: %s\n", signedCookie)
	}

	responseData := map[string]interface{}{
		"user": authResponse.User,
	}

	successResponse := common.NewSuccessResponse(responseData)
	common.WriteJSONResponse(w, successResponse, 200)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := common.NewErrorResponse(400, "Solicitud inválida")
		common.WriteJSONResponse(w, response, 400)
		return
	}

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

	var currentUser *auth.User
	if user, ok := ctx.Value("user").(*auth.User); ok {
		currentUser = user
	}

	authResponse, signedCookie, expiresAt, err := h.authService.RegisterUser(ctx, req, currentUser)
	if err != nil {
		statusCode := 400
		if strings.Contains(err.Error(), "Ha ocurrido un error en el servidor") {
			statusCode = 500
		} else if strings.Contains(err.Error(), "no tiene permisos") {
			statusCode = 403
		}

		errorResponse := common.NewErrorResponse(statusCode, err.Error())
		common.WriteJSONResponse(w, errorResponse, statusCode)
		return
	}

	if signedCookie != "" {
		cookies.SetAuthCookie(w, signedCookie, expiresAt)
		fmt.Printf("🎯 COOKIE SET IN REGISTER: %s\n", signedCookie)
	}

	responseData := map[string]interface{}{
		"user": authResponse.User,
	}

	successResponse := common.NewSuccessResponse(responseData)
	common.WriteJSONResponse(w, successResponse, 201)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if userIDStr, ok := ctx.Value("userID").(string); ok {
		if userID, err := strconv.Atoi(userIDStr); err == nil {
			h.authService.Logout(ctx, userID)
		}
	}

	cookies.ClearAuthCookie(w)

	successResponse := common.NewSuccessResponse("Sesión cerrada exitosamente")
	common.WriteJSONResponse(w, successResponse, 200)
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*auth.User)
	if !ok || user == nil {
		response := common.NewErrorResponse(401, "No autenticado")
		common.WriteJSONResponse(w, response, 401)
		return
	}

	safeUser := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"rol":      user.RolNombre,
	}

	successResponse := common.NewSuccessResponse(safeUser)
	common.WriteJSONResponse(w, successResponse, 200)
}
