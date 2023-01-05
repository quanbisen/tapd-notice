package cron

import (
	gocron "github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
	"tapd-notice/common/dingding"
	"tapd-notice/model"
)

type DingdingDeptCronJob struct {
	Orm    *gorm.DB
	Client dingding.Client
}

func NewDingdingDeptCronJob(orm *gorm.DB, client dingding.Client) gocron.Job {
	return &DingdingDeptCronJob{
		Orm:    orm,
		Client: client,
	}
}

func (c *DingdingDeptCronJob) Run() {
	log.Println("DingdingDeptCronJob Run...")
	depts, err := c.Client.ListDept()
	if err != nil {
		log.Printf("DingdingDeptCronJob ListDept failed, err: %s", err)
		return
	}
	for i := range depts {
		dept := model.DingdingDept{}
		if err = c.Orm.Model(&dept).Find(&dept, "dept_id=?", depts[i].DeptId).Error; err != nil {
			log.Printf("DingdingDeptCronJob Find DingdingDept failed, dept_id=%d, err: %s", depts[i].DeptId, err)
			continue
		}
		newDept := model.DingdingDept{
			DeptId:          depts[i].DeptId,
			ParentId:        depts[i].ParentId,
			Name:            depts[i].Name,
			AutoAddUser:     depts[i].AutoAddUser,
			CreateDeptGroup: depts[i].CreateDeptGroup,
		}
		// 不存在就创建
		if dept.DeptId == 0 {
			if err = c.Orm.Model(&dept).Create(newDept).Error; err != nil {
				log.Printf("DingdingDeptCronJob Create DingdingDept failed, dept: %v", dept)
			}
		} else {
			if err = c.Orm.Model(&dept).Updates(&newDept).Error; err != nil {
				log.Printf("DingdingDeptCronJob Update DingdingDept failed, dept: %v", dept)
			}
		}
	}
	log.Println("DingdingDeptCronJob End...")
}
