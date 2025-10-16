package whatsapp

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, controller *WhatsAppController) {
	whatsappGroup := router.Group("/whatsapp")
	{
		sessionsGroup := whatsappGroup.Group("/sessions")
		{
			sessionsGroup.POST("", controller.RegisterSession)
			sessionsGroup.GET("", controller.ListSessions)
			sessionsGroup.GET("/:sessionId", controller.GetSessionStatus)
			sessionsGroup.DELETE("/:sessionId", controller.DeleteSession)
		}

		whatsappGroup.POST("/messages", controller.SendMessage)
	}
}
