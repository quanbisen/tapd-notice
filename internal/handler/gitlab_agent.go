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

type GitlabAgentApi struct {
	api.Api
}

func (api GitlabAgentApi) SendMessage(c *gin.Context) {
	agentMessageData := dto.GitlabAgentMessageData{}
	s := service.NewGitlabAgentService(config.GetDingdingConfig().GitlabAgent)
	err := api.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&agentMessageData, binding.JSON).Errors
	if err != nil {
		log.Printf("Api prepare failed, err: %v\n", err)
		api.Error(http.StatusInternalServerError, err, "Api prepare error")
		return
	}
	err = s.SendMessage(&agentMessageData)
	if err != nil {
		log.Printf("GitlabAgentApi SendMessage failed, err: %v\n", err)
		api.Error(http.StatusInternalServerError, err, "GitlabAgentApi SendMessage error")
		return
	}
	log.Printf("GitlabAgentApi SendMessage success\n")
	api.OK(nil, "GitlabAgentApi SendMessage success")
}
