package orm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"tapd-notice/config"
)

var DB *gorm.DB

func Setup(opts ...func(db *gorm.DB)) {
	databaseConfig := config.GetDatabaseConfig()
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: databaseConfig.Source}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("failed to connect database")
	}
	resolver := &dbresolver.DBResolver{}
	resolver.SetMaxOpenConns(databaseConfig.MaxOpenConns)
	resolver.SetMaxIdleConns(databaseConfig.MaxIdleConns)
	db.Use(resolver)
	DB = db.Debug()
	for _, opt := range opts {
		opt(DB)
	}
}
