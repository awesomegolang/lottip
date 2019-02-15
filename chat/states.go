package chat

// Cmd represents MySQL command to be executed.
type Cmd struct {
	ConnId     uint
	CmdId      uint
	Executable bool
	Database   string
	Query      string
	Parameters []string
}

// CmdResult represents MySQL command execution result.
type CmdResult struct {
	ConnId   uint
	CmdId    uint
	Result   byte
	Error    string
	Duration string
}

// ConnState represents tcp connection state.
type ConnState struct {
	ConnId uint
	State  byte
}
