package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"tapd-notice/common/api"
	"tapd-notice/config"
	"tapd-notice/internal/dto"
	"tapd-notice/internal/service"
)

type WebhookApi struct {
	api.Api
}

func (api *WebhookApi) SendWebhook(c *gin.Context) {
	webhook := dto.WebhookData{}
	s := service.NewWebhookService(config.GetDingdingConfig())
	err := api.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&webhook, binding.JSON).Errors
	if err != nil {
		api.Error(http.StatusInternalServerError, err, "Api prepare error")
		return
	}
	s.PrintRawData()
	err = s.SendMessage(&webhook)
	if err != nil {
		api.Error(http.StatusInternalServerError, err, "SendMessage error")
		return
	}
	api.OK(nil, "SendMessage success")
}
