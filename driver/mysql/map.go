package mysql

import "time"

// Command represents MySQL command.
type Command struct {
	ConnId     uint          `json:"conn_id"`
	CmdId      uint          `json:"cmd_id"`
	Result     byte          `json:"result"`
	Error      string        `json:"error"`
	Query      string        `json:"query"`
	Database   string        `json:"database"`
	Executable bool          `json:"executable"`
	Parameters []string      `json:"parameters"`
	Duration   time.Duration `json:"duration"`
}

// Connection represents tcp connection.
type Connection struct {
	ConnId uint `json:"conn_id"`
	State  byte `json:"state"`
}
