package cron

import (
	gocron "github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
	"tapd-notice/common/tapd"
	"tapd-notice/model"
	"time"
)

type TAPDProjectCronJob struct {
	Orm    *gorm.DB
	Client tapd.Client
}

func NewTAPDProjectCronJob(orm *gorm.DB, client tapd.Client) gocron.Job {
	return &TAPDProjectCronJob{
		Orm:    orm,
		Client: client,
	}
}

func (c *TAPDProjectCronJob) Run() {
	log.Println("TAPDProjectCronJob Run...")
	projects, err := c.Client.ListProject()
	if err != nil {
		log.Printf("TAPDProjectCronJob ListProject failed, err: %s", err)
		return
	}
	duration, _ := time.ParseDuration("-8h")
	localLocation, _ := time.LoadLocation("Local")
	for i := range projects {
		project := model.TAPDProject{}
		if err = c.Orm.Model(&project).Find(&project, "id=?", projects[i].Id).Error; err != nil {
			log.Printf("TAPDProjectCronJob Find TAPDProject failed, id=%s, err: %s", projects[i].Id, err)
			continue
		}

		newProject := model.TAPDProject{
			Id:          projects[i].Id,
			Name:        projects[i].Name,
			Description: projects[i].Description,
			Status:      projects[i].Status,
			ParentId:    projects[i].ParentId,
			Created:     time.Time(projects[i].Created).Add(duration).In(localLocation),
			Creator:     projects[i].Creator,
			MemberCount: projects[i].MemberCount,
		}
		// 不存在就创建
		if project.Id == "" {
			if err = c.Orm.Model(&project).Create(newProject).Error; err != nil {
				log.Printf("TAPDProjectCronJob Create TAPDProject failed, project: %v", project)
			}
		} else {
			if err = c.Orm.Model(&project).Updates(&newProject).Error; err != nil {
				log.Printf("TAPDProjectCronJob Update TAPDProject failed, project: %v", project)
			}
		}
	}
	log.Println("TAPDProjectCronJob End...")
}
