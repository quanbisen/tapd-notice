package migrate

import (
	"gorm.io/gorm"
	"log"
	"tapd-notice/model"
)

func WithMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&model.TAPDProject{}, &model.DingdingDept{}, &model.DingdingUser{})
	if err != nil {
		log.Fatalf("migrate model err, error: %s", err)
	}
}
