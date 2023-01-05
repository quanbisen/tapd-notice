package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tapd-notice/internal/handler"
	"tapd-notice/middleware"
)

func NewRouter() *gin.Engine {
	engine := gin.New()
	engine.UseH2C = true
	engine.Use(gin.Logger(),
		gin.Recovery(),
		middleware.NoCache(),
		middleware.Options(),
		middleware.Secure(),
		middleware.ContextDB())

	engine.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	v1 := engine.Group("/api/v1")
	{
		webhookApi := handler.WebhookApi{}
		v1.POST("/webhook", webhookApi.SendWebhook)
	}
	return engine
}
