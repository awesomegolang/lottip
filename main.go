package main

import (
	"encoding/json"
	"github.com/orderbynull/lottip/driver/mysql"
	"github.com/orderbynull/lottip/driver/redis"
)

func main() {
	var ready chan interface{}

	for _, driver := range drivers() {
		switch cfg := driver.(type) {

		case MySQLConfig:
			proxy := mysql.NewProxyServer(cfg.Target, cfg.Listen)
			go proxy.Run()
			go func() {
				for {
					select {
					case cmd := <-proxy.Commands:
						jsonData, _ := json.Marshal(cmd)
						println("MySQL -> ", string(jsonData))
						//go postData(apiUrl(apiHost, mysqlCmdPath), jsonData)

					case conn := <-proxy.Connections:
						jsonData, _ := json.Marshal(conn)
						println("CONN -> ", string(jsonData))
						//go postData(apiUrl(mysqlConnPath), jsonData)
					}
				}
			}()
			
		case RedisConfig:
			monitor := redis.NewMonitor(cfg.Target)
			go monitor.Run()
			go func() {
				for cmd := range monitor.Commands {
					println("Redis: ", cmd.Query)
				}
			}()
		}
	}

	<- ready
}
