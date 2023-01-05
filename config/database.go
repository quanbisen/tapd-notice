package config

import "github.com/spf13/viper"

type DatabaseConfig struct {
	Source       string
	MaxIdleConns int
	MaxOpenConns int
}

func initDatabase(v *viper.Viper) DatabaseConfig {
	databaseConfig := DatabaseConfig{
		Source:       v.GetString("source"),
		MaxIdleConns: v.GetInt("maxidleconns"),
		MaxOpenConns: v.GetInt("maxopenconns"),
	}
	if databaseConfig.MaxOpenConns == 0 {
		databaseConfig.MaxOpenConns = 5
	}
	if databaseConfig.MaxIdleConns == 0 {
		databaseConfig.MaxIdleConns = 20
	}
	validateDatabase(databaseConfig)
	return databaseConfig
}

func validateDatabase(config DatabaseConfig) {
	if config.Source == "" {
		panic("database.source config empty.")
	}
}
