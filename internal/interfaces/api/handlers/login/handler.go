package login

import (
	"encoding/json"
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

// LoginRequestDTO específico para el handler
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
	if r.Method != http.MethodPost {
		response := common.NewErrorResponse(405, "Método no permitido")
		common.WriteJSONResponse(w, response, 405)
		return
	}

	var reqDTO LoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		response := common.NewErrorResponse(400, "Solicitud JSON inválida")
		common.WriteJSONResponse(w, response, 400)
		return
	}

	if err := reqDTO.Validate(); err != nil {
		utils.WriteValidationError(w, err, 400)
		return
	}

	ctx := r.Context()
	credentials := login.LoginCredentials{
		Username: strings.TrimSpace(reqDTO.Username),
		Password: reqDTO.Password,
	}

	response, err := h.loginService.Execute(ctx, credentials)
	if err != nil {
		errorResponse := common.NewErrorResponse(401, err.Error())
		common.WriteJSONResponse(w, errorResponse, 401)
		return
	}

	common.WriteJSONResponse(w, *response, 200)
}
