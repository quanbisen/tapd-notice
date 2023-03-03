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

func (s *WebhookService) PrintRawData() error {
	data := make(map[string]interface{}, 0)
	if v, exist := s.C.Get("body"); exist {
		if byts, ok := v.([]byte); ok {
			if err := json.Unmarshal(byts, &data); err != nil {
				log.Printf("WebhookService PrintRawData failed, json.Unmarshal error: %s\n", err)
				return fmt.Errorf("WebhookService PrintRawData failed, json.Unmarshal error: %s\n", err)
			}
			newByts, err := json.Marshal(data)
			if err != nil {
				log.Printf("WebhookService PrintRawData failed, json.Marshal error: %s\n", err)
				return fmt.Errorf("WebhookService PrintRawData failed, json.Marshal error: %s\n", err)
			}
			log.Println("WebhookService received message start...")
			log.Printf("\n%s\n", string(newByts))
			log.Println("WebhookService received message end...")
		}
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
		res += fmt.Sprintf("**创建人：** %s  \n  ", strings.Join(data.Event.Creator, ","))
		res += fmt.Sprintf("**当前状态：** %s  \n  ", common.StoryStatusMap[data.Event.Status])
		res += fmt.Sprintf("**处理人：** %s  ", strings.Join(data.Event.Owner, ","))
	} else if data.Event.EventKey == common.StoryStatusChange {
		accessUrl := fmt.Sprintf("https://www.tapd.cn/%s/prong/stories/view/%s", data.WorkspaceId, data.Event.Id)
		res += fmt.Sprintf("### %s：%s->%s  \n  ", common.EventKeyMap[common.StoryStatusChange], common.StoryStatusMap[data.Event.StatusFromTo.From], common.StoryStatusMap[data.Event.StatusFromTo.To])
		res += fmt.Sprintf("**项目：** %s  \n  ", projectName)
		res += fmt.Sprintf("**标题：** [%s](%s)  \n  ", data.Event.Name, accessUrl)
		res += fmt.Sprintf("**创建人：** %s  \n  ", strings.Join(data.Event.Creator, ","))
		res += fmt.Sprintf("**当前状态：** %s  \n  ", common.StoryStatusMap[data.Event.Status])
		res += fmt.Sprintf("**处理人：** %s  ", strings.Join(data.Event.Owner, ","))
	} else if data.Event.EventKey == common.StoryComment {
		accessUrl := fmt.Sprintf("https://www.tapd.cn/%s/prong/stories/view/%s", data.WorkspaceId, data.Event.Id)
		res += fmt.Sprintf("### %s：%s  \n  ", common.EventKeyMap[common.StoryComment], common.CommentActionMap[data.Event.SubEvent])
		res += fmt.Sprintf("**项目：** %s  \n  ", projectName)
		res += fmt.Sprintf("**标题：** [%s](%s)  \n  ", data.Event.Name, accessUrl)
		if data.Event.SubEvent == common.CommentAdd || data.Event.SubEvent == common.CommentUpdate {
			res += fmt.Sprintf("**评论内容：** %s  \n  ", util.ExtractComment(data.Event.DescriptionFromTo.To))
		} else if data.Event.SubEvent == common.CommentDelete {
			res += fmt.Sprintf("**原评论内容：** *%s*  \n  ", util.ExtractComment(data.Event.DescriptionFromTo.To))
		}
	} else if data.Event.EventKey == common.BugCreate {
		accessUrl := fmt.Sprintf("https://www.tapd.cn/%s/bugtrace/bugs/view?bug_id=%s", data.WorkspaceId, data.Event.Id)
		res += fmt.Sprintf("### %s  \n  ", common.EventKeyMap[common.BugCreate])
		res += fmt.Sprintf("**项目：** %s  \n  ", projectName)
		res += fmt.Sprintf("**标题：** [%s](%s)  \n  ", data.Event.Title, accessUrl)
		res += fmt.Sprintf("**创建人：** %s  \n  ", strings.Join(data.Event.Reporter, ","))
		res += fmt.Sprintf("**当前状态：** %s  \n  ", common.BugStatusMap[data.Event.Status])
		res += fmt.Sprintf("**处理人：** %s  ", strings.Join(data.Event.CurrentOwner, ","))
	} else if data.Event.EventKey == common.BugStatusChange {
		accessUrl := fmt.Sprintf("https://www.tapd.cn/%s/bugtrace/bugs/view?bug_id=%s", data.WorkspaceId, data.Event.Id)
		res += fmt.Sprintf("### %s：%s->%s  \n  ", common.EventKeyMap[common.BugStatusChange], common.BugStatusMap[data.Event.StatusFromTo.From], common.BugStatusMap[data.Event.StatusFromTo.To])
		res += fmt.Sprintf("**项目：** %s  \n  ", projectName)
		res += fmt.Sprintf("**标题：** [%s](%s)  \n  ", data.Event.Title, accessUrl)
		res += fmt.Sprintf("**创建人：** %s  \n  ", strings.Join(data.Event.Reporter, ","))
		res += fmt.Sprintf("**当前状态：** %s  \n  ", common.BugStatusMap[data.Event.Status])
		res += fmt.Sprintf("**处理人：** %s  ", strings.Join(data.Event.CurrentOwner, ","))
	} else if data.Event.EventKey == common.BugComment {
		accessUrl := fmt.Sprintf("https://www.tapd.cn/%s/bugtrace/bugs/view?bug_id=%s", data.WorkspaceId, data.Event.Id)
		res += fmt.Sprintf("### %s：%s  \n  ", common.EventKeyMap[common.BugComment], common.CommentActionMap[data.Event.SubEvent])
		res += fmt.Sprintf("**项目：** %s  \n  ", projectName)
		res += fmt.Sprintf("**标题：** [%s](%s)  \n  ", data.Event.Title, accessUrl)
		if data.Event.SubEvent == common.CommentAdd || data.Event.SubEvent == common.CommentUpdate {
			res += fmt.Sprintf("**评论内容：** %s  \n  ", util.ExtractComment(data.Event.DescriptionFromTo.To))
		} else if data.Event.SubEvent == common.CommentDelete {
			res += fmt.Sprintf("**原评论内容：** *%s*  \n  ", util.ExtractComment(data.Event.DescriptionFromTo.To))
		}
	} else if data.Event.EventKey == common.TaskCreate {
		accessUrl := fmt.Sprintf("https://www.tapd.cn/%s/prong/tasks/view/%s", data.WorkspaceId, data.Event.Id)
		res += fmt.Sprintf("### %s  \n  ", common.EventKeyMap[common.TaskCreate])
		res += fmt.Sprintf("**项目：** %s  \n  ", projectName)
		res += fmt.Sprintf("**标题：** [%s](%s)  \n  ", data.Event.Name, accessUrl)
		res += fmt.Sprintf("**创建人：** %s  \n  ", strings.Join(data.Event.Creator, ","))
		res += fmt.Sprintf("**当前状态：** %s  \n  ", common.TaskStatusMap[data.Event.Status])
		res += fmt.Sprintf("**处理人：** %s  ", strings.Join(data.Event.Owner, ","))
	} else if data.Event.EventKey == common.TaskStatusChange {
		accessUrl := fmt.Sprintf("https://www.tapd.cn/%s/prong/tasks/view/%s", data.WorkspaceId, data.Event.Id)
		res += fmt.Sprintf("### %s：%s->%s  \n  ", common.EventKeyMap[common.TaskStatusChange], common.TaskStatusMap[data.Event.StatusFromTo.From], common.TaskStatusMap[data.Event.StatusFromTo.To])
		res += fmt.Sprintf("**项目：** %s  \n  ", projectName)
		res += fmt.Sprintf("**标题：** [%s](%s)  \n  ", data.Event.Name, accessUrl)
		res += fmt.Sprintf("**创建人：** %s  \n  ", strings.Join(data.Event.Creator, ","))
		res += fmt.Sprintf("**当前状态：** %s  \n  ", common.TaskStatusMap[data.Event.Status])
		res += fmt.Sprintf("**处理人：** %s  ", strings.Join(data.Event.Owner, ","))
	} else if data.Event.EventKey == common.TaskComment {
		accessUrl := fmt.Sprintf("https://www.tapd.cn/%s/prong/tasks/view/%s", data.WorkspaceId, data.Event.Id)
		res += fmt.Sprintf("### %s：%s  \n  ", common.EventKeyMap[common.TaskComment], common.CommentActionMap[data.Event.SubEvent])
		res += fmt.Sprintf("**项目：** %s  \n  ", projectName)
		res += fmt.Sprintf("**标题：** [%s](%s)  \n  ", data.Event.Name, accessUrl)
		if data.Event.SubEvent == common.CommentAdd || data.Event.SubEvent == common.CommentUpdate {
			res += fmt.Sprintf("**评论内容：** %s  \n  ", util.ExtractComment(data.Event.DescriptionFromTo.To))
		} else if data.Event.SubEvent == common.CommentDelete {
			res += fmt.Sprintf("**原评论内容：** *%s*  \n  ", util.ExtractComment(data.Event.DescriptionFromTo.To))
		}
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
		adviseUsers = append(adviseUsers, data.Event.Creator...)
		adviseUsers = append(adviseUsers, data.Event.Owner...)
	} else if data.Event.EventKey == common.StoryComment {
		// 针对评论处理逻辑，如果评论是@人的，单独通知被@的人，否则，通知创建人和处理人
		cuePeople := util.ExtractCommentCuePeople(data.Event.DescriptionFromTo.To)
		if len(cuePeople) > 0 {
			adviseUsers = append(adviseUsers, cuePeople...)
		} else {
			adviseUsers = append(adviseUsers, data.Event.Creator...)
			adviseUsers = append(adviseUsers, data.Event.Owner...)
		}
	} else if data.Event.EventKey == common.BugCreate {
		adviseUsers = append(adviseUsers, data.Event.CurrentOwner...)
	} else if data.Event.EventKey == common.BugStatusChange {
		adviseUsers = append(adviseUsers, data.Event.Reporter...)
		adviseUsers = append(adviseUsers, data.Event.CurrentOwner...)
	} else if data.Event.EventKey == common.BugComment {
		// 针对评论处理逻辑，如果评论是@人的，单独通知被@的人，否则，通知创建人和处理人
		cuePeople := util.ExtractCommentCuePeople(data.Event.DescriptionFromTo.To)
		if len(cuePeople) > 0 {
			adviseUsers = append(adviseUsers, cuePeople...)
		} else {
			adviseUsers = append(adviseUsers, data.Event.Reporter...)
			adviseUsers = append(adviseUsers, data.Event.CurrentOwner...)
		}
	} else if data.Event.EventKey == common.TaskCreate {
		adviseUsers = append(adviseUsers, data.Event.Owner...)
		adviseUsers = append(adviseUsers, data.Event.CC...)
	} else if data.Event.EventKey == common.TaskStatusChange {
		adviseUsers = append(adviseUsers, data.Event.Creator...)
		adviseUsers = append(adviseUsers, data.Event.Owner...)
		adviseUsers = append(adviseUsers, data.Event.CC...)
	} else if data.Event.EventKey == common.TaskComment {
		// 针对评论处理逻辑，如果评论是@人的，单独通知被@的人，否则，通知创建人和处理人
		cuePeople := util.ExtractCommentCuePeople(data.Event.DescriptionFromTo.To)
		if len(cuePeople) > 0 {
			adviseUsers = append(adviseUsers, cuePeople...)
		} else {
			adviseUsers = append(adviseUsers, data.Event.Creator...)
			adviseUsers = append(adviseUsers, data.Event.Owner...)
			adviseUsers = append(adviseUsers, data.Event.CC...)
		}
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
	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": title,
			"text":  text,
		},
	}
	return s.dingdingClient.SendAppMessage(msg, userIds)
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
