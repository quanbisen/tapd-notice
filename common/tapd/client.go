package tapd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"tapd-notice/internal/dto"
)

type Client struct {
	client      *http.Client
	companyId   string
	apiUser     string
	apiPassword string
}

func NewClient(companyId, apiUser, apiPassword string) Client {
	client := Client{
		client:      &http.Client{},
		companyId:   companyId,
		apiUser:     apiUser,
		apiPassword: apiPassword,
	}

	return client
}

func (c *Client) ListProject() ([]dto.TAPDProject, error) {
	url := c.getUrlPrefix() + fmt.Sprintf("/workspaces/projects?company_id=%s", c.companyId)
	byts, err := c.doRequest(url, "GET", nil)
	if err != nil {
		log.Printf("TAPD Client ListProject failed, doRequest err: %s\n", err)
		return nil, err
	}
	var listRes dto.TAPDProjectListResult
	if err = json.Unmarshal(byts, &listRes); err != nil {
		log.Printf("TAPD Client ListProject failed, json.Unmarshal error: %s\n", err)
		return nil, err
	}
	if listRes.Status == 1 && listRes.Info == "success" {
		res := make([]dto.TAPDProject, 0)
		for _, workspace := range listRes.Data {
			project := workspace.Workspace
			res = append(res, project)
		}
		return res, nil
	} else {
		return nil, fmt.Errorf("TAPD Client ListProject failed, message: %s", listRes.Info)
	}

}

func (c *Client) doRequest(url string, method string, data interface{}) ([]byte, error) {
	var (
		jsonBytes []byte
		err       error
		request   *http.Request
	)

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
	basicAuthBase64 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.apiUser, c.apiPassword)))
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuthBase64))
	res, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed at status code: %d", res.StatusCode)
	}
	return io.ReadAll(res.Body)
}

func (c *Client) getUrlPrefix() string {
	return "https://api.tapd.cn"
}
