package model

import "tapd-notice/common"

type DingdingUser struct {
	UserId   string `json:"userid" gorm:"primaryKey;types:varchar(128)"`
	Name     string `json:"name" gorm:"types:varchar(64)"`
	Username string `json:"username" gorm:"types:varchar(128);index"`
	DeptId   int64  `json:"dept_id" gorm:"types:int(11);column:dept_id"`
	Email    string `json:"email" gorm:"types:varchar(64)"`
	Title    string `json:"title" gorm:"types:varchar(64)"`
	Mobile   string `json:"mobile" gorm:"types:varchar(32)"`
}

func (user *DingdingUser) TableName() string {
	return common.DingdingUserTableName
}
