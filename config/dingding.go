package config

import "github.com/spf13/viper"

type DingdingConfig struct {
	AppKey    string
	AppSecret string
	AgentId   int64
}

func initDingding(v *viper.Viper) DingdingConfig {
	dingdingConfig := DingdingConfig{
		AppKey:    v.GetString("appkey"),
		AppSecret: v.GetString("appsecret"),
		AgentId:   v.GetInt64("agentid"),
	}
	validateDingding(dingdingConfig)
	return dingdingConfig
}

func validateDingding(config DingdingConfig) {
	if config.AppKey == "" {
		panic("dingding.appkey config empty.")
	}
	if config.AppSecret == "" {
		panic("dingding.appsecret config empty.")
	}
	if config.AgentId == 0 {
		panic("dingding.agentid config empty.")
	}
}
