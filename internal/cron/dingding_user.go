package cron

import (
	"github.com/mozillazg/go-pinyin"
	gocron "github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
	"strings"
	"tapd-notice/common/dingding"
	"tapd-notice/model"
)

type DingdingUserCronJob struct {
	Orm    *gorm.DB
	Client dingding.Client
}

func NewDingdingUserCronJob(orm *gorm.DB, client dingding.Client) gocron.Job {
	return &DingdingUserCronJob{
		Orm:    orm,
		Client: client,
	}
}

func (c *DingdingUserCronJob) Run() {
	log.Println("DingdingUserCronJob Run...")
	db := c.Orm.Begin()
	var (
		depts []model.DingdingDept
		users []model.DingdingUser
	)
	if err := db.Model(&model.DingdingDept{}).Find(&depts).Error; err != nil {
		log.Printf("DingdingUserCronJob Find depts failed, err: %s", err)
		return
	}
	args := pinyin.NewArgs()
	// 遍历获取部门下的所有用户
	for _, dept := range depts {
		deptUsers, err := c.Client.ListDeptUser(dept.DeptId)
		if err != nil {
			log.Printf("DingdingUserCronJob ListDeptUser failed, dept_id: %d", dept.DeptId)
			continue
		}
		for i := range deptUsers {
			user := model.DingdingUser{
				UserId:   deptUsers[i].UserId,
				Name:     deptUsers[i].Name,
				Username: strings.Join(pinyin.LazyPinyin(deptUsers[i].Name, args), ""),
				DeptId:   dept.DeptId,
			}
			users = append(users, user)
		}
	}
	for i := range users {
		user := model.DingdingUser{}
		if err := db.Model(&user).Find(&user, "user_id=? and dept_id=?", users[i].UserId, users[i].DeptId).Error; err != nil {
			log.Printf("DingdingUserCronJob Find DingdingUser failed, user_id=%s, dept_id=%d, err: %s", users[i].UserId, depts[i].DeptId, err)
			continue
		}
		newUser := users[i]
		// 不存在就创建
		if user.UserId == "" {
			if err := db.Model(&user).Create(newUser).Error; err != nil {
				log.Printf("DingdingUserCronJob Create DingdingUser failed, user: %v", newUser)
			}
		} else {
			if err := db.Model(&user).Updates(&newUser).Error; err != nil {
				log.Printf("DingdingUserCronJob Update DingdingUser failed, user: %v", newUser)
			}
		}
	}
	db.Commit()
	log.Println("DingdingUserCronJob End...")
}
