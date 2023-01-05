package dto

type FromTo struct {
	To   string `json:"to"`
	From string `json:"from"`
}

type Event struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Title        string   `json:"title"`
	User         string   `json:"user"`
	EventKey     string   `json:"event_key"`
	Owner        []string `json:"owner"`
	CurrentOwner []string `json:"current_owner"`
	CC           []string `json:"cc"`
	Status       string   `json:"status"`
	StatusFromTo FromTo   `json:"status:fromto"`
}

type WebhookData struct {
	WorkspaceId string `json:"workspace_id"`
	Event       Event  `json:"event"`
}

type WebhookResult struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
