package dingding

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestAuthenticate(t *testing.T) {
	appKey := os.Getenv("appkey")
	appSecret := os.Getenv("appsecret")
	client := NewClient(appKey, appSecret, "")
	err := client.Authenticate()
	assert.Empty(t, err)
	assert.NotEmpty(t, client.accessToken)
}

func TestDoRequestAccessTokenExpired(t *testing.T) {
	appKey := os.Getenv("appkey")
	appSecret := os.Getenv("appsecret")
	client := NewClient(appKey, appSecret, "")
	sub, _ := time.ParseDuration("-4h")
	client.expiredAt = client.expiredAt.Add(sub)
	user, err := client.ListDeptUser(54213524)
	assert.Empty(t, err)
	assert.NotEmpty(t, user)
}

func TestListDept(t *testing.T) {
	appKey := os.Getenv("appkey")
	appSecret := os.Getenv("appsecret")
	client := NewClient(appKey, appSecret, "")
	depts, err := client.ListDept()
	assert.Empty(t, err)
	assert.NotEmpty(t, depts)
}

func TestListDeptUser(t *testing.T) {
	appKey := os.Getenv("appkey")
	appSecret := os.Getenv("appsecret")
	client := NewClient(appKey, appSecret, "")
	user, err := client.ListDeptUser(54213524)
	assert.Empty(t, err)
	assert.NotEmpty(t, user)
}

func TestSendAppMessage(t *testing.T) {
	appKey := os.Getenv("appkey")
	appSecret := os.Getenv("appsecret")
	agentId := os.Getenv("agentid")
	client := NewClient(appKey, appSecret, agentId)
	res, err := client.SendAppMessage("需求创建", "### 需求创建  \n  **标题：** [测试需求](https://www.tapd.cn/50628422/prong/stories/view/1150628422001003390)  \n  **创建人：** 全碧森  \n  **当前状态：** 规划中  \n  **处理人：** quanbisen,韦炳铁  ", "16601997543784086")
	assert.Empty(t, err)
	assert.Equal(t, "ok", res.ErrMsg)
}
