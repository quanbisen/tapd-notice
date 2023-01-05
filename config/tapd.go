package config

import "github.com/spf13/viper"

type TAPDConfig struct {
	CompanyId   string
	ApiUser     string
	ApiPassword string
}

func initTAPD(v *viper.Viper) TAPDConfig {
	tapdConfig := TAPDConfig{
		CompanyId:   v.GetString("companyid"),
		ApiUser:     v.GetString("apiuser"),
		ApiPassword: v.GetString("apipassword"),
	}
	validateTAPD(tapdConfig)
	return tapdConfig
}

func validateTAPD(config TAPDConfig) {
	if config.CompanyId == "" {
		panic("tapd.companyid config empty.")
	}
	if config.ApiUser == "" {
		panic("tapd.apiuser config empty.")
	}
	if config.ApiPassword == "" {
		panic("tapd.apipassword config empty.")
	}
}
