package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
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
		log.Printf("Api prepare failed, err: %v\n", err)
		api.Error(http.StatusInternalServerError, err, "Api prepare error")
		return
	}
	s.PrintRawData()
	err = s.SendMessage(&webhook)
	if err != nil {
		log.Printf("SendMessage failed, err: %v\n", err)
		api.Error(http.StatusInternalServerError, err, "SendMessage error")
		return
	}
	log.Printf("SendMessage success\n")
	api.OK(nil, "SendMessage success")
}
