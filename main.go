package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/orderbynull/lottip/chat"
	"time"
)

const mysqlAddr = "127.0.0.1:3306"
const proxyAddr = "127.0.0.1:4041"



func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&MySql{})

	cmdChan := make(chan chat.Cmd)
	cmdResultChan := make(chan chat.CmdResult)
	connStateChan := make(chan chat.ConnState)
	appReadyChan := make(chan bool)

	go func() {
		for {
			select {
			case cmd := <-cmdChan:
				db.Create(&MySql{CmdId: cmd.CmdId, ConnId: cmd.ConnId, Query: cmd.Query, Done: false, Duration: time.Second})
			case result := <-cmdResultChan:
				db.Model(&MySql{}).Where("cmd_id = ? AND conn_id = ?", result.CmdId, result.ConnId).Updates(map[string]interface{}{"error": result.Error, "done": true})
			case <-connStateChan:
			}
		}
	}()

	(&MySQLProxyServer{cmdChan, cmdResultChan, connStateChan, appReadyChan, mysqlAddr, proxyAddr}).run()
}
