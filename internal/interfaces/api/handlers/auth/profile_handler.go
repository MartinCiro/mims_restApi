package auth

import (
	"log"
	"net/http"

	"api_go/internal/core/auth"
	"api_go/internal/interfaces/api/common"
)

type ProfileHandler struct {
	authService *auth.AuthService
}

func NewProfileHandler(authService *auth.AuthService) *ProfileHandler {
	return &ProfileHandler{
		authService: authService,
	}
}

// GetProfile compatible con http.HandlerFunc
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// ✅ Agregar recover para capturar panics
	defer func() {
		if err := recover(); err != nil {
			log.Printf("❌ PANIC en GetProfile: %v", err)
			common.InternalServerError(w, "Error interno del servidor")
		}
	}()

	log.Printf("🎯 Iniciando GetProfile para userID del contexto")

	// Obtener userID del contexto (seteado por el middleware de auth)
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		log.Printf("❌ No se pudo obtener userID del contexto. ok: %v, userID: '%v'", ok, userID)
		common.Unauthorized(w, "Usuario no autenticado")
		return
	}

	log.Printf("✅ userID obtenido del contexto: %s", userID)

	// Obtener perfil del usuario usando el servicio
	response, err := h.authService.GetUserProfile(r.Context(), userID)
	if err != nil {
		log.Printf("❌ Error en GetUserProfile: %v", err)
		// Manejar diferentes tipos de error
		switch err.Error() {
		case "ID de usuario inválido":
			common.BadRequest(w, err.Error())
		case "usuario no encontrado":
			common.NotFound(w, err.Error())
		default:
			common.InternalServerError(w, "Error al obtener perfil: "+err.Error())
		}
		return
	}

	log.Printf("✅ Perfil obtenido exitosamente, enviando respuesta")

	// ✅ Validar que response no sea nil
	if response == nil {
		log.Printf("❌ Response es nil")
		common.InternalServerError(w, "Respuesta inválida del servicio")
		return
	}

	// Devolver la respuesta
	common.WriteJSONResponse(w, response, response.Code)
}
