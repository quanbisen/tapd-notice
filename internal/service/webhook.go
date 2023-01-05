package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tapd-notice/common"
	"tapd-notice/common/dingding"
	"tapd-notice/common/util"
	"tapd-notice/config"
	"tapd-notice/internal/dto"
	"tapd-notice/model"
)

type WebhookService struct {
	client         http.Client
	dingdingClient dingding.Client
	Service
}

func NewWebhookService(dingdingConfig config.DingdingConfig) *WebhookService {
	dingdingClient := dingding.NewClient(dingdingConfig.AppKey, dingdingConfig.AppSecret, strconv.FormatInt(dingdingConfig.AgentId, 10))
	return &WebhookService{
		dingdingClient: dingdingClient,
		client:         http.Client{},
	}
}

func (s *WebhookService) SendMessage(data *dto.WebhookData) error {
	message := s.generateMarkdownMessage(data)
	res, err := s.sendMessage(common.EventKeyMap[data.Event.EventKey], message, data)
	if err != nil {
		log.Printf("WebhookService SendMessage failed,  err: %s\n", err)
		return fmt.Errorf("WebhookService SendMessage failed,  err: %s\n", err)
	}
	if res.ErrMsg != "ok" {
		log.Printf("WebhookService SendMessage failed, msg: %s\n", res.ErrMsg)
		return fmt.Errorf("WebhookService SendMessage failed, msg: %s\n", res.ErrMsg)
	}
	return nil
}

func (s *WebhookService) generateMarkdownMessage(data *dto.WebhookData) string {
	var (
		res         string
		projectName string
	)
	project := model.TAPDProject{Id: data.WorkspaceId}
	if err := s.Orm.Model(&project).Find(&project).Error; err != nil {
		log.Printf("WebhookService generateMarkdownMessage Find TAPDProject failed, err: %s\n", err)
	}
	projectName = project.Name

	if data.Event.EventKey == common.StoryCreate {
		accessUrl := fmt.Sprintf("https://www.tapd.cn/%s/prong/stories/view/%s", data.WorkspaceId, data.Event.Id)
		res += fmt.Sprintf("### %s  \n  ", common.EventKeyMap[common.StoryCreate])
		res += fmt.Sprintf("**项目：** %s  \n  ", projectName)
		res += fmt.Sprintf("**标题：** [%s](%s)  \n  ", data.Event.Name, accessUrl)
		res += fmt.Sprintf("**创建人：** %s  \n  ", data.Event.User)
		res += fmt.Sprintf("**当前状态：** %s  \n  ", common.StoryStatusMap[data.Event.Status])
		res += fmt.Sprintf("**处理人：** %s  ", strings.Join(data.Event.Owner, ","))
	} else if data.Event.EventKey == common.StoryStatusChange {
		accessUrl := fmt.Sprintf("https://www.tapd.cn/%s/prong/stories/view/%s", data.WorkspaceId, data.Event.Id)
		res += fmt.Sprintf("### %s：%s->%s  \n  ", common.EventKeyMap[common.StoryStatusChange], common.StoryStatusMap[data.Event.StatusFromTo.From], common.StoryStatusMap[data.Event.StatusFromTo.To])
		res += fmt.Sprintf("**项目：** %s  \n  ", projectName)
		res += fmt.Sprintf("**标题：** [%s](%s)  \n  ", data.Event.Name, accessUrl)
		res += fmt.Sprintf("**创建人：** %s  \n  ", data.Event.User)
		res += fmt.Sprintf("**当前状态：** %s  \n  ", common.StoryStatusMap[data.Event.Status])
		res += fmt.Sprintf("**处理人：** %s  ", strings.Join(data.Event.Owner, ","))
	} else if data.Event.EventKey == common.BugCreate {
		accessUrl := fmt.Sprintf("https://www.tapd.cn/%s/bugtrace/bugs/view?bug_id=%s", data.WorkspaceId, data.Event.Id)
		res += fmt.Sprintf("### %s  \n  ", common.EventKeyMap[common.BugCreate])
		res += fmt.Sprintf("**项目：** %s  \n  ", projectName)
		res += fmt.Sprintf("**标题：** [%s](%s)  \n  ", data.Event.Title, accessUrl)
		res += fmt.Sprintf("**创建人：** %s  \n  ", data.Event.User)
		res += fmt.Sprintf("**当前状态：** %s  \n  ", common.BugStatusMap[data.Event.Status])
		res += fmt.Sprintf("**处理人：** %s  ", strings.Join(data.Event.CurrentOwner, ","))
	} else if data.Event.EventKey == common.BugStatusChange {
		accessUrl := fmt.Sprintf("https://www.tapd.cn/%s/bugtrace/bugs/view?bug_id=%s", data.WorkspaceId, data.Event.Id)
		res += fmt.Sprintf("### %s：%s->%s  \n  ", common.EventKeyMap[common.BugStatusChange], common.BugStatusMap[data.Event.StatusFromTo.From], common.BugStatusMap[data.Event.StatusFromTo.To])
		res += fmt.Sprintf("**项目：** %s  \n  ", projectName)
		res += fmt.Sprintf("**标题：** [%s](%s)  \n  ", data.Event.Title, accessUrl)
		res += fmt.Sprintf("**创建人：** %s  \n  ", data.Event.User)
		res += fmt.Sprintf("**当前状态：** %s  \n  ", common.BugStatusMap[data.Event.Status])
		res += fmt.Sprintf("**处理人：** %s  ", strings.Join(data.Event.CurrentOwner, ","))
	} else if data.Event.EventKey == common.TaskCreate {
		accessUrl := fmt.Sprintf("https://www.tapd.cn/%s/prong/tasks/view/%s", data.WorkspaceId, data.Event.Id)
		res += fmt.Sprintf("### %s  \n  ", common.EventKeyMap[common.TaskCreate])
		res += fmt.Sprintf("**项目：** %s  \n  ", projectName)
		res += fmt.Sprintf("**标题：** [%s](%s)  \n  ", data.Event.Name, accessUrl)
		res += fmt.Sprintf("**创建人：** %s  \n  ", data.Event.User)
		res += fmt.Sprintf("**当前状态：** %s  \n  ", common.TaskStatusMap[data.Event.Status])
		res += fmt.Sprintf("**处理人：** %s  ", strings.Join(data.Event.Owner, ","))
	} else if data.Event.EventKey == common.TaskStatusChange {
		accessUrl := fmt.Sprintf("https://www.tapd.cn/%s/prong/tasks/view/%s", data.WorkspaceId, data.Event.Id)
		res += fmt.Sprintf("### %s：%s->%s  \n  ", common.EventKeyMap[common.TaskStatusChange], common.TaskStatusMap[data.Event.StatusFromTo.From], common.TaskStatusMap[data.Event.StatusFromTo.To])
		res += fmt.Sprintf("**项目：** %s  \n  ", projectName)
		res += fmt.Sprintf("**标题：** [%s](%s)  \n  ", data.Event.Name, accessUrl)
		res += fmt.Sprintf("**创建人：** %s  \n  ", data.Event.User)
		res += fmt.Sprintf("**当前状态：** %s  \n  ", common.TaskStatusMap[data.Event.Status])
		res += fmt.Sprintf("**处理人：** %s  ", strings.Join(data.Event.Owner, ","))
	}
	return res
}

func (s *WebhookService) sendMarkdownMessageToRobot(robotURL string, title string, text string) (*dto.WebhookResult, error) {
	data := make(map[string]interface{}, 0)
	data["msgtype"] = "markdown"
	data["markdown"] = map[string]string{"title": title, "text": text}
	return s.doRequest(robotURL, data)
}

func (s *WebhookService) sendMessage(title string, text string, data *dto.WebhookData) (*dto.DingdingSendAppMessageResult, error) {
	userIds := make([]string, 0)
	adviseUsers := make([]string, 0)
	if data.Event.EventKey == common.StoryCreate {
		adviseUsers = append(adviseUsers, data.Event.Owner...)
	} else if data.Event.EventKey == common.StoryStatusChange {
		adviseUsers = append(adviseUsers, data.Event.User)
		adviseUsers = append(adviseUsers, data.Event.Owner...)
	} else if data.Event.EventKey == common.BugCreate {
		adviseUsers = append(adviseUsers, data.Event.CurrentOwner...)
	} else if data.Event.EventKey == common.BugStatusChange {
		adviseUsers = append(adviseUsers, data.Event.User)
		adviseUsers = append(adviseUsers, data.Event.CurrentOwner...)
	} else if data.Event.EventKey == common.TaskCreate {
		adviseUsers = append(adviseUsers, data.Event.Owner...)
		adviseUsers = append(adviseUsers, data.Event.CC...)
	} else if data.Event.EventKey == common.TaskStatusChange {
		adviseUsers = append(adviseUsers, data.Event.User)
		adviseUsers = append(adviseUsers, data.Event.Owner...)
		adviseUsers = append(adviseUsers, data.Event.CC...)
	}
	for _, adviseUser := range adviseUsers {
		username := util.CharToPinyin(adviseUser)
		user := model.DingdingUser{}
		if err := s.Orm.Model(&user).First(&user, "username=?", username).Error; err != nil {
			log.Printf("WebhookService sendMessage failed, Find user error: %s\n", err)
			continue
		}
		userIds = append(userIds, user.UserId)
	}
	userIds = util.Unique(userIds)
	return s.dingdingClient.SendAppMessage(title, text, userIds)
}

func (s *WebhookService) doRequest(url string, data interface{}) (*dto.WebhookResult, error) {
	var (
		jsonBytes []byte
		err       error
	)
	jsonBytes, err = json.Marshal(data)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	res, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	result := dto.WebhookResult{}
	if err = json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
