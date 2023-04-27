package model

import "tapd-notice/common"

type DingdingUser struct {
	UserId   string `json:"userid" gorm:"primaryKey;type:varchar(128)"`
	Name     string `json:"name" gorm:"type:varchar(64)"`
	Username string `json:"username" gorm:"type:varchar(128);index"`
	DeptId   int64  `json:"dept_id" gorm:"type:int(11);column:dept_id"`
	Email    string `json:"email" gorm:"type:varchar(64)"`
	Title    string `json:"title" gorm:"type:varchar(64)"`
	Mobile   string `json:"mobile" gorm:"type:varchar(32)"`
}

func (user *DingdingUser) TableName() string {
	return common.DingdingUserTableName
}
