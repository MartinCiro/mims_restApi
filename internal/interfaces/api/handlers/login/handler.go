// internal/interfaces/api/handlers/login/handler.go
package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"api_go/internal/core/login"
	"api_go/internal/interfaces/api/common"
	"api_go/pkg/utils"
)

type LoginHandler struct {
	loginService *login.LoginService
}

func NewLoginHandler(loginService *login.LoginService) *LoginHandler {
	return &LoginHandler{
		loginService: loginService,
	}
}

type LoginRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (dto *LoginRequestDTO) Validate() error {
	if strings.TrimSpace(dto.Username) == "" {
		return utils.ValidarBlank("username", "El nombre de usuario es obligatorio")
	}
	if strings.TrimSpace(dto.Password) == "" {
		return utils.ValidarBlank("password", "La contraseña es obligatoria")
	}
	return nil
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("🔍 Headers: %v\n", r.Header)
	fmt.Printf("🔍 Content-Type: %s\n", r.Header.Get("Content-Type"))

	// Leer el body completo para debug
	bodyBytes, _ := io.ReadAll(r.Body)
	fmt.Printf("🔍 Raw Body: %s\n", string(bodyBytes))

	// Resetear el body
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var reqDTO LoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		fmt.Printf("❌ Error decoding JSON: %v\n", err)
		response := common.NewErrorResponse(400, "Solicitud JSON inválida: "+err.Error())
		common.WriteJSONResponse(w, response, 400)
		return
	}

	fmt.Printf("🔍 Request DTO después de decode: Username='%s', Password length=%d\n",
		reqDTO.Username,
		len(reqDTO.Password))

	// Validar campos
	if err := reqDTO.Validate(); err != nil {
		fmt.Printf("❌ Validación falló: %v\n", err)
		utils.WriteValidationError(w, err, 400)
		return
	}

	fmt.Printf("✅ JSON válido, procediendo con login...\n")

	ctx := r.Context()
	credentials := login.LoginCredentials{
		Username: strings.TrimSpace(reqDTO.Username),
		Password: reqDTO.Password,
	}

	// Ejecutar servicio del core (retorna objeto del dominio)
	result, err := h.loginService.Execute(ctx, credentials)
	if err != nil {
		fmt.Printf("❌ Error en login service: %v\n", err)
		errorResponse := common.NewErrorResponse(401, err.Error())
		common.WriteJSONResponse(w, errorResponse, 401)
		return
	}

	// Adaptar resultado del dominio a respuesta HTTP
	fmt.Printf("✅ Login exitoso para: %s\n", reqDTO.Username)
	response := common.NewSuccessResponse(result)
	common.WriteJSONResponse(w, response, 200)
}
