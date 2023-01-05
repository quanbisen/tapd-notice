package dto

import (
	"tapd-notice/common/types"
)

type TAPDProject struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	PrettyName  string          `json:"pretty_name"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	ParentId    string          `json:"parent_id"`
	Secrecy     string          `json:"secrecy"`
	Created     types.LocalTime `json:"created"`
	CreatorId   string          `json:"creator_id"`
	Creator     string          `json:"creator"`
	MemberCount int             `json:"member_count"`
}

type TAPDProjectListResult struct {
	Status int `json:"status"`
	Data   []struct {
		Workspace TAPDProject `json:"Workspace"`
	} `json:"data"`
	Info string `json:"info"`
}
