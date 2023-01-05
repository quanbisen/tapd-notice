package config

import (
	"github.com/spf13/viper"
	"sync"
)

var once sync.Once
var settings Settings

type Settings struct {
	ApplicationConfig ApplicationConfig
	DingdingConfig    DingdingConfig
	TAPDConfig        TAPDConfig
	DatabaseConfig    DatabaseConfig
}

func Init(host, port, configYml string) {
	once.Do(func() {
		viper.SetConfigFile(configYml)
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		// application
		application := viper.Sub("settings.application")
		applicationConfig := initApplication(host, port, application)
		// dingding
		dingding := viper.Sub("settings.dingding")
		dingdingConfig := initDingding(dingding)
		// tapd
		tapd := viper.Sub("settings.tapd")
		tapdConfig := initTAPD(tapd)
		// database
		database := viper.Sub("settings.database")
		databaseConfig := initDatabase(database)

		settings = Settings{
			ApplicationConfig: applicationConfig,
			DingdingConfig:    dingdingConfig,
			TAPDConfig:        tapdConfig,
			DatabaseConfig:    databaseConfig,
		}
	})
}

func GetApplicationConfig() ApplicationConfig {
	return settings.ApplicationConfig
}

func GetDingdingConfig() DingdingConfig {
	return settings.DingdingConfig
}

func GetTAPDConfig() TAPDConfig {
	return settings.TAPDConfig
}

func GetDatabaseConfig() DatabaseConfig {
	return settings.DatabaseConfig
}
