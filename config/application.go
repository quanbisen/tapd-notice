package config

import "github.com/spf13/viper"

type ApplicationConfig struct {
	Host string
	Port string
}

func initApplication(host, port string, v *viper.Viper) ApplicationConfig {
	applicationConfig := ApplicationConfig{
		Host: v.GetString("host"),
		Port: v.GetString("port"),
	}
	if host != "" {
		applicationConfig.Host = host
	}
	if port != "" {
		applicationConfig.Port = port
	}
	validateApplication(applicationConfig)
	return applicationConfig
}

func validateApplication(config ApplicationConfig) {
	if config.Host == "" {
		panic("application.host config empty.")
	}
	if config.Port == "" {
		panic("application.port config empty.")
	}
}
