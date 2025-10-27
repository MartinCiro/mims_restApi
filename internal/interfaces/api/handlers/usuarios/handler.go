package usuarios

import (
	"net/http"

	"api_go/internal/interfaces/api/common"
)

type UsuariosHandler struct{}

func NewUsuariosHandler() *UsuariosHandler {
	return &UsuariosHandler{}
}

func (h *UsuariosHandler) ObtenerUsuarios(w http.ResponseWriter, r *http.Request) {
	response := common.NewSuccessResponse(map[string]string{
		"message": "Endpoint de usuarios - por implementar",
	})
	common.WriteJSONResponse(w, response, 200)
}
