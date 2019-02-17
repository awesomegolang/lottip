package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/orderbynull/lottip/driver/mysql"
	"net/http"
)

const mysqlAddr = "127.0.0.1:3306"
const proxyAddr = "127.0.0.1:4041"
const apiUrl = "https://enqjpyw74ueq.x.pipedream.net"

func mysqlRoute() string {
	return fmt.Sprintf("%s", apiUrl)
}

func postData(apiUrl string, data []byte) {
	_, _ = http.Post(apiUrl, "application/json", bytes.NewBuffer(data))
}

func main() {
	proxy := mysql.NewProxyServer(mysqlAddr, proxyAddr)

	go func() {
		for {
			select {
			case cmd := <-proxy.Commands:
				jsonData, _ := json.Marshal(cmd)
				go postData(mysqlRoute(), jsonData)
			case <-proxy.Connections:
			}
		}
	}()

	proxy.Run()
}
