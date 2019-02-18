package main

import (
	"encoding/json"
	"github.com/orderbynull/lottip/driver/mysql"
)

const mysqlAddr = "127.0.0.1:3306"
const proxyAddr = "127.0.0.1:4041"
const apiHost = "https://enqjpyw74ueq.x.pipedream.net"
const mysqlCmdPath = "/mysql/cmd"
const mysqlConnPath = "/mysql/conn"

func main() {
	proxy := mysql.NewProxyServer(mysqlAddr, proxyAddr)

	go func() {
		for {
			select {
			case cmd := <-proxy.Commands:
				jsonData, _ := json.Marshal(cmd)
				println("CMD -> ", string(jsonData))
				//go postData(apiUrl(mysqlCmdPath), jsonData)

			case <-proxy.Connections:
				//jsonData, _ := json.Marshal(conn)
				//println("CONN -> ", string(jsonData))
				//go postData(apiUrl(mysqlConnPath), jsonData)
			}
		}
	}()

	proxy.Run()
}
