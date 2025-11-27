package auth

import (
	"encoding/json"
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
		common.WriteSimpleError(w, "Solicitud inválida", 400)
		return
	}

	if req.Email == "" {
		common.WriteSimpleError(w, "El email es requerido", 400)
		return
	}

	if req.Password == "" {
		common.WriteSimpleError(w, "La contraseña es requerida", 400)
		return
	}

	ctx := r.Context()

	loginReq := auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	authResponse, signedCookie, expiresAt, err := h.authService.LoginUser(ctx, loginReq)
	if err != nil {
		common.WriteSimpleError(w, err.Error(), 401)
		return
	}

	if signedCookie != "" {
		cookies.SetAuthCookie(w, signedCookie, expiresAt)
	}

	successResponse := common.NewSuccessResponse(authResponse.Message)
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
	}

	// CAMBIO AQUÍ: Manejar diferentes tipos de respuesta
	var responseData interface{}

	if currentUser != nil {
		responseData = authResponse
	} else {
		// Para usuarios no autenticados (registro normal)
		// Extraer el mensaje directamente
		switch v := authResponse.(type) {
		case *auth.AuthResponse:
			responseData = v.Message
		case string:
			responseData = v
		default:
			responseData = "Usuario registrado con exito"
		}
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
