package config

import "github.com/spf13/viper"

type Dingding struct {
	TapdAgent   DingdingConfig
	GitlabAgent DingdingConfig
}

type DingdingConfig struct {
	AppKey    string
	AppSecret string
	AgentId   int64
}

func initDingding(v *viper.Viper) Dingding {
	dingding := Dingding{
		TapdAgent: DingdingConfig{
			AppKey:    v.GetString("tapdagent.appkey"),
			AppSecret: v.GetString("tapdagent.appsecret"),
			AgentId:   v.GetInt64("tapdagent.agentid"),
		},
		GitlabAgent: DingdingConfig{
			AppKey:    v.GetString("gitlabagent.appkey"),
			AppSecret: v.GetString("gitlabagent.appsecret"),
			AgentId:   v.GetInt64("gitlabagent.agentid"),
		},
	}
	validateDingdingConfig(dingding.TapdAgent)
	validateDingdingConfig(dingding.GitlabAgent)
	return dingding
}

func validateDingdingConfig(config DingdingConfig) {
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
