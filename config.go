package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type MySQLConfig struct {
	Listen string
	Target string
}

type RedisConfig struct {
	Target string
}

func drivers() []interface{} {
	var drivers []interface{}
	drivers = append(
		drivers,
		MySQLConfig{viper.GetString("drivers.mysql.listen"), viper.GetString("drivers.mysql.target")},
		RedisConfig{viper.GetString("drivers.redis.target")},
	)

	return drivers
}

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
