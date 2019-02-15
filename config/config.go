package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	configName = "config"
	configType = "yaml"
	configPath = "."

	appDbHostTmpl = "applications.%s.%s.host"
	appDbPortTmpl = "applications.%s.%s.port"
	appProxyHostTmpl = "applications.%s.proxy.host"
	appProxyPortTmpl = "applications.%s.proxy.port"
)

type Application struct {
	Name      string
	Driver    string
	DbPort    int
	DbHost    string
	ProxyPort int
	ProxyHost string
}

type Config struct {
	Applications []Application
}

func Read() *Config {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	var driver string
	var applications []Application

	for appName, _ := range viper.GetStringMapString("applications") {
		driver = viper.GetString(fmt.Sprintf("applications.%s.driver", appName))
		app := Application{
			Name:      appName,
			Driver:    driver,
			DbHost:    viper.GetString(fmt.Sprintf(appDbHostTmpl, appName, driver)),
			DbPort:    viper.GetInt(fmt.Sprintf(appDbPortTmpl, appName, driver)),
			ProxyHost: viper.GetString(fmt.Sprintf(appProxyHostTmpl, appName)),
			ProxyPort: viper.GetInt(fmt.Sprintf(appProxyPortTmpl, appName)),
		}

		applications = append(applications, app)
	}

	return &Config{Applications:applications}
}
