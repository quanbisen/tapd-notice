package service

import (
	"fmt"
	"gorm.io/gorm"
)

type Service struct {
	Orm   *gorm.DB
	Error error
}

func (db *Service) AddError(err error) error {
	if db.Error == nil {
		db.Error = err
	} else if err != nil {
		db.Error = fmt.Errorf("%v; %w", db.Error, err)
	}
	return db.Error
}
