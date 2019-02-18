package main

import (
	"encoding/json"
	"fmt"
	"github.com/orderbynull/lottip/driver/mysql"
	"github.com/spf13/viper"
)

const mysqlCmdPath = "/mysql/cmd"
const mysqlConnPath = "/mysql/conn"

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	var proxyAddr = viper.GetString("proxies.mysql.in")
	var mysqlAddr = viper.GetString("proxies.mysql.out")
	//var apiHost = viper.GetString("api_url")
	
	proxy := mysql.NewProxyServer(mysqlAddr, proxyAddr)

	go func() {
		for {
			select {
			case cmd := <-proxy.Commands:
				jsonData, _ := json.Marshal(cmd)
				println("CMD -> ", string(jsonData))
				//go postData(apiUrl(apiHost, mysqlCmdPath), jsonData)

			case conn := <-proxy.Connections:
				jsonData, _ := json.Marshal(conn)
				println("CONN -> ", string(jsonData))
				//go postData(apiUrl(mysqlConnPath), jsonData)
			}
		}
	}()

	proxy.Run()
}
