package service

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tapd-notice/common/dingding"
	"tapd-notice/config"
	"tapd-notice/internal/dto"
	"tapd-notice/model"
)

type AgentMessageService struct {
	client         http.Client
	dingdingClient dingding.Client
	Service
}

func NewGitlabAgentService(dingdingConfig config.DingdingConfig) *AgentMessageService {
	dingdingClient := dingding.NewClient(dingdingConfig.AppKey, dingdingConfig.AppSecret, strconv.FormatInt(dingdingConfig.AgentId, 10))
	return &AgentMessageService{
		dingdingClient: dingdingClient,
		client:         http.Client{},
	}
}

func (s *AgentMessageService) SendMessage(messageData *dto.GitlabAgentMessageData) error {
	user := model.DingdingUser{}
	err := s.Orm.Model(&user).Where("username=?", messageData.Username).First(&user).Error
	if err != nil {
		log.Printf("AgentMessageService Find user failed, err: %s\n", err)
		return err
	}
	res, err := s.dingdingClient.SendAppMessage(messageData.Msg, []string{user.UserId})
	if err != nil {
		log.Printf("AgentMessageService SendMessage failed,  err: %s\n", err)
		return fmt.Errorf("AgentMessageService SendMessage failed,  err: %s\n", err)
	}
	if res.ErrMsg != "ok" {
		log.Printf("AgentMessageService SendMessage failed, msg: %s\n", res.ErrMsg)
		return fmt.Errorf("AgentMessageService SendMessage failed, msg: %s\n", res.ErrMsg)
	}
	return nil
}
