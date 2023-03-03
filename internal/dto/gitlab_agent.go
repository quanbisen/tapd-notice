package dto

type GitlabAgentMessageData struct {
	Username string                 `json:"username" binding:"required"`
	Msg      map[string]interface{} `json:"msg" binding:"required"`
}
