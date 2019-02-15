package model

import "time"

type MySql struct {
	CmdId     uint `gorm:"primary_key, auto_increment:false"`
	ConnId    uint `gorm:"primary_key, auto_increment:false"`
	Done      bool
	Query     string
	Error     string
	CreatedAt time.Time
	Duration  time.Duration `gorm:"type:int"`
}

func (MySql) TableName() string {
	return "mysql"
}
