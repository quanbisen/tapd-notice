package dingding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"tapd-notice/common/util"
	"tapd-notice/internal/dto"
	"time"
)

type Client struct {
	client      *http.Client
	accessToken string
	expiredAt   time.Time
	appKey      string
	appSecret   string
	agentId     string
}

func NewClient(appKey, appSecret, agentId string) Client {
	client := Client{
		client:    &http.Client{},
		appKey:    appKey,
		appSecret: appSecret,
		agentId:   agentId,
	}

	return client
}

func (c *Client) Authenticate() error {
	url := c.getUrlPrefix() + fmt.Sprintf("/gettoken?appkey=%s&appsecret=%s", c.appKey, c.appSecret)
	byts, err := c.doRequest(url, "GET", nil)
	if err != nil {
		log.Printf("Dingding Client Authenticate failed, doRequest error: %s\n", err)
		return err
	}
	res := dto.DingdingGetTokenResult{}
	if err = json.Unmarshal(byts, &res); err != nil {
		log.Printf("Dingding Client Authenticate failed, json.Unmarshal error: %s\n", err)
		return err
	}
	c.accessToken = res.AccessToken
	c.expiredAt = time.Now().Add(time.Second * time.Duration(res.ExpiresIn))
	return nil
}

func (c *Client) ListDept() ([]dto.DingdingDept, error) {
	res := make([]dto.DingdingDept, 0)
	leafNode, err := c.listDept(nil)
	if err != nil {
		log.Printf("Dingding Client listDept failed, error: %s\n", err)
		return nil, err
	}
	res = append(res, leafNode...)
	for len(leafNode) > 0 {
		rangeDepts := leafNode
		leafNode = leafNode[0:0]
		for i := 0; i < len(rangeDepts); i++ {
			depts, err := c.listDept(util.Int64Addr(rangeDepts[i].DeptId))
			if err != nil {
				log.Printf("Dingding Client listDept failed, error: %s\n", err)
				return nil, err
			}
			res = append(res, depts...)
			leafNode = append(leafNode, depts...)
		}
	}
	return res, nil
}

func (c *Client) listDept(deptId *int64) ([]dto.DingdingDept, error) {
	var res dto.DingdingDeptListResult
	url := c.getUrlPrefix() + fmt.Sprintf("/topapi/v2/department/listsub?access_token=%s", c.accessToken)
	var data map[string]int64
	if deptId != nil {
		data = map[string]int64{"dept_id": *deptId}
	}
	byts, err := c.doRequest(url, "POST", data)
	if err != nil {
		log.Printf("Dingding Client listDept failed, error: %s\n", err)
		return nil, err
	}
	if err = json.Unmarshal(byts, &res); err != nil {
		log.Printf("Dingding Client listDept failed, json.Unmarshal error: %s\n", err)
		return nil, err
	}
	if res.ErrCode != 0 || res.ErrMsg != "ok" {
		log.Printf("Dingding Client listDept failed, error: %s\n", res.ErrMsg)
		return nil, fmt.Errorf("Dingding Client listDept failed, error: %s", res.ErrMsg)
	}
	return res.Result, nil
}

func (c *Client) ListDeptUser(deptId int64) ([]dto.DingdingDeptUser, error) {
	var (
		res     []dto.DingdingDeptUser
		listRes *dto.DingdingDeptUserListResult
	)
	data := make(map[string]interface{}, 3)
	data["dept_id"] = deptId
	data["cursor"] = 0
	data["size"] = 1
	listRes, err := c.listDeptUser(data)
	if err != nil {
		log.Printf("Dingding Client ListDeptUser failed, error: %s\n", err)
		return nil, err
	}
	res = append(res, listRes.Result.List...)
	for listRes.Result.HasMore {
		cursor, _ := util.GetIntByInterface(data["cursor"])
		data["cursor"] = cursor + 1
		listRes, err = c.listDeptUser(data)
		if err != nil {
			log.Printf("Dingding Client ListDeptUser failed, error: %s\n", err)
			return nil, err
		}
		res = append(res, listRes.Result.List...)
	}
	return res, nil
}

func (c *Client) listDeptUser(data map[string]interface{}) (*dto.DingdingDeptUserListResult, error) {
	var res dto.DingdingDeptUserListResult
	url := c.getUrlPrefix() + fmt.Sprintf("/topapi/user/listsimple?access_token=%s", c.accessToken)
	byts, err := c.doRequest(url, "POST", data)
	if err != nil {
		log.Printf("Dingding Client listDeptUser failed, error: %s\n", err)
		return nil, err
	}
	if err = json.Unmarshal(byts, &res); err != nil {
		log.Printf("Dingding Client listDeptUser failed, json.Unmarshal error: %s\n", err)
		return nil, err
	}
	if res.ErrCode != 0 || res.ErrMsg != "ok" {
		log.Printf("Dingding Client listDeptUser failed, error: %s\n", res.ErrMsg)
		return nil, fmt.Errorf("Dingding Client listDeptUser failed, error: %s", res.ErrMsg)
	}
	return &res, nil
}

func (c *Client) SendAppMessage(title, text string, userIdList []string) (*dto.DingdingSendAppMessageResult, error) {
	if len(userIdList) == 0 {
		return nil, nil
	}
	data := make(map[string]interface{}, 0)
	data["agent_id"] = c.agentId
	data["userid_list"] = strings.Join(userIdList, ",")
	data["to_all_user"] = false
	data["msg"] = map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": title,
			"text":  text,
		},
	}
	var res dto.DingdingSendAppMessageResult
	url := c.getUrlPrefix() + fmt.Sprintf("/topapi/message/corpconversation/asyncsend_v2?access_token=%s", c.accessToken)
	byts, err := c.doRequest(url, "POST", data)
	if err != nil {
		log.Printf("Dingding Client SendAppMessage failed, doRequest error: %s\n", err)
		return nil, err
	}
	if err = json.Unmarshal(byts, &res); err != nil {
		log.Printf("Dingding Client SendAppMessage failed, json.Unmarshal error: %s\n", err)
		return nil, err
	}
	return &res, nil
}

func (c *Client) doRequest(url string, method string, data interface{}) ([]byte, error) {
	var (
		jsonBytes []byte
		err       error
		request   *http.Request
	)
	// 验证access_token是否已经过期
	if !strings.Contains(url, "gettoken") {
		nowTime := time.Now()
		if nowTime.After(c.expiredAt) {
			if err = c.Authenticate(); err != nil {
				log.Printf("Dingding Client Authenticate failed, error: %s\n", err)
				return nil, err
			}
			if strings.Contains(url, "access_token=") {
				// 替换url上的access_token字段值
				regex, _ := regexp.Compile("access_token=(\\w{32}|\\w{0})")
				url = regex.ReplaceAllString(url, fmt.Sprintf("access_token=%s", c.accessToken))
			}
			c.doRequest(url, method, data)
		}
	}
	if data != nil {
		jsonBytes, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
		request, err = http.NewRequest(method, url, bytes.NewReader(jsonBytes))
		if err != nil {
			return nil, err
		}
	} else {
		request, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	res, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(res.Body)
}

func (c *Client) getUrlPrefix() string {
	return "https://oapi.dingtalk.com"
}
