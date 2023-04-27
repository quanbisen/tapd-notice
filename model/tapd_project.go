package model

import (
	"tapd-notice/common"
	"time"
)

type TAPDProject struct {
	Id          string    `json:"id" gorm:"primaryKey;type:varchar(16)"`
	Name        string    `json:"name" gorm:"type:varchar(256)"`
	Description string    `json:"description" gorm:"type:varchar(512)"`
	Status      string    `json:"status" gorm:"type:varchar(16)"`
	ParentId    string    `json:"parent_id" gorm:"column:parent_id;type:varchar(16)"`
	Created     time.Time `json:"created"`
	Creator     string    `json:"creator" gorm:"type:varchar(128)"`
	MemberCount int       `json:"member_count" gorm:"column:member_count"`
}

func (p *TAPDProject) TableName() string {
	return common.TAPDProjectTableName
}
