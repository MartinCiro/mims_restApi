package whatsapp

import (
	"scrapper_go_email/core/whatsapp"
	"scrapper_go_email/interfaces/api/common"
	"time"

	"github.com/gin-gonic/gin"
)

type WhatsAppController struct {
	service whatsapp.WhatsAppService
}

func NewWhatsAppController(service whatsapp.WhatsAppService) *WhatsAppController {
	return &WhatsAppController{
		service: service,
	}
}

// RegisterSession godoc
// @Summary Registrar nueva sesión de WhatsApp
// @Description Crea una nueva sesión y retorna el QR code para autenticación
// @Tags whatsapp
// @Accept json
// @Produce json
// @Param body body RegisterRequestDTO true "Datos de registro"
// @Success 201 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Router /api/v1/whatsapp/sessions [post]
func (c *WhatsAppController) RegisterSession(ctx *gin.Context) {
	var req RegisterRequestDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Datos inválidos: "+err.Error())
		return
	}

	// Convertir DTO a dominio
	registerReq := &whatsapp.RegisterRequest{
		SessionID:  req.SessionID,
		Provider:   req.Provider,
		WebhookURL: req.WebhookURL,
		Metadata:   req.Metadata,
	}

	result, err := c.service.RegisterSession(ctx.Request.Context(), registerReq)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}

	// Crear respuesta con los datos del registro
	response := SessionResponseDTO{
		Provider:   req.Provider,
		Status:     string(result.Status),
		QRCode:     result.QRCode,
		WebhookURL: req.WebhookURL,
		Metadata:   req.Metadata,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	common.Created(ctx, response)
}

// GetSessionStatus godoc
// @Summary Obtener estado de sesión
// @Description Retorna el estado actual de una sesión de WhatsApp
// @Tags whatsapp
// @Produce json
// @Param sessionId path string true "ID de la sesión"
// @Success 200 {object} common.SuccessResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /api/v1/whatsapp/sessions/{sessionId} [get]
func (c *WhatsAppController) GetSessionStatus(ctx *gin.Context) {
	sessionID := ctx.Param("sessionId")

	session, err := c.service.GetSessionStatus(ctx.Request.Context(), sessionID)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}

	// Convertir Session del core a DTO
	sessionDTO := convertSessionToDTO(session)
	common.Success(ctx, sessionDTO)
}

// SendMessage godoc
// @Summary Enviar mensaje
// @Description Envía un mensaje a través de WhatsApp
// @Tags whatsapp
// @Accept json
// @Produce json
// @Param body body MessageRequestDTO true "Datos del mensaje"
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Router /api/v1/whatsapp/messages [post]
func (c *WhatsAppController) SendMessage(ctx *gin.Context) {
	var req MessageRequestDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Datos inválidos: "+err.Error())
		return
	}

	messageReq := &whatsapp.MessageRequest{
		SessionID: req.SessionID,
		To:        req.To,
		Message:   req.Message,
		Type:      req.Type,
	}

	result, err := c.service.SendMessage(ctx.Request.Context(), messageReq)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}

	// Usar SuccessWithData para formato { data: ... }
	common.Success(ctx, result.Data)
}

// DeleteSession godoc
// @Summary Eliminar sesión
// @Description Elimina una sesión de WhatsApp
// @Tags whatsapp
// @Param sessionId path string true "ID de la sesión"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /api/v1/whatsapp/sessions/{sessionId} [delete]
func (c *WhatsAppController) DeleteSession(ctx *gin.Context) {
	sessionID := ctx.Param("sessionId")

	if err := c.service.DeleteSession(ctx.Request.Context(), sessionID); err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.NoContent(ctx)
}

// ListSessions godoc
// @Summary Listar sesiones
// @Description Retorna todas las sesiones registradas
// @Tags whatsapp
// @Produce json
// @Success 200 {object} common.SuccessResponse
// @Router /api/v1/whatsapp/sessions [get]
func (c *WhatsAppController) ListSessions(ctx *gin.Context) {
	sessions, err := c.service.ListSessions(ctx.Request.Context())
	if err != nil {
		common.HandleError(ctx, err)
		return
	}

	// Convertir sessions a DTOs
	sessionsDTO := make([]SessionResponseDTO, len(sessions))
	for i, session := range sessions {
		sessionsDTO[i] = convertSessionToDTO(session)
	}

	common.Success(ctx, sessionsDTO)
}

// Función helper para convertir Session del core a DTO
func convertSessionToDTO(session *whatsapp.Session) SessionResponseDTO {
	dto := SessionResponseDTO{
		ID:         session.ID,
		Provider:   session.Provider,
		Status:     string(session.Status),
		QRCode:     session.QRCode,
		WebhookURL: session.WebhookURL,
		Metadata:   session.Metadata,
		CreatedAt:  session.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  session.UpdatedAt.Format(time.RFC3339),
	}

	if session.ConnectedAt != nil {
		dto.ConnectedAt = session.ConnectedAt.Format(time.RFC3339)
	}

	if session.UserInfo != nil {
		dto.UserInfo = &UserInfoDTO{
			Phone:      session.UserInfo.Phone,
			Name:       session.UserInfo.Name,
			ProfilePic: session.UserInfo.ProfilePic,
			Status:     session.UserInfo.Status,
		}
	}

	return dto
}
