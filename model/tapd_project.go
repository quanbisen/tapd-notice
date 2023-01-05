package model

import (
	"tapd-notice/common"
	"time"
)

type TAPDProject struct {
	Id          string    `json:"id" gorm:"primaryKey;types:varchar(16)"`
	Name        string    `json:"name" gorm:"types:varchar(256)"`
	Description string    `json:"description" gorm:"types:varchar(512)"`
	Status      string    `json:"status" gorm:"types:varchar(16)"`
	ParentId    string    `json:"parent_id" gorm:"column:parent_id;types:varchar(16)"`
	Created     time.Time `json:"created"`
	Creator     string    `json:"creator" gorm:"types:varchar(128)"`
	MemberCount int       `json:"member_count" gorm:"column:member_count"`
}

func (p *TAPDProject) TableName() string {
	return common.TAPDProjectTableName
}
