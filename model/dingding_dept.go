package model

import "tapd-notice/common"

type DingdingDept struct {
	DeptId          int64  `json:"dept_id" gorm:"primaryKey;type:int(11);column:dept_id"`
	ParentId        int64  `json:"parent_id" gorm:"type:int(11);column:parent_id"`
	Name            string `json:"name" gorm:"type:varchar(128)"`
	AutoAddUser     bool   `json:"auto_add_user" gorm:"column:auto_add_user"`
	CreateDeptGroup bool   `json:"create_dept_group" gorm:"column:create_dept_group"`
}

func (dept *DingdingDept) TableName() string {
	return common.DingdingDeptTableName
}
