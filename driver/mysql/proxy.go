package mysql

import (
	"io"
	"log"
	"net"
	"time"
)

type RequestPacketParser struct {
	connId      uint
	queryId     *uint
	timer       *time.Time
	commandsMap map[uint]*Command
}

func (pp *RequestPacketParser) Write(p []byte) (n int, err error) {
	*pp.queryId++
	*pp.timer = time.Now()

	switch packetType(p) {
	case comStmtPrepare:
	case comQuery:
		decoded, err := decodeQueryRequest(p)
		if err == nil {
			pp.commandsMap[*pp.queryId] = &Command{
				ConnId:     pp.connId,
				CmdId:      *pp.queryId,
				Query:      decoded.Query,
				Executable: false,
			}
		}
	case comQuit:
		println("CLOSED")
		//pp.connStateChan <- Connection{pp.connId, connStateFinished}
	}

	return len(p), nil
}

type ResponsePacketParser struct {
	connId    uint
	cmdId     *uint
	queryChan chan Command
	timer     *time.Time
	toRename  map[uint]*Command
}

func (pp *ResponsePacketParser) Write(p []byte) (n int, err error) {
	if command, ok := pp.toRename[*pp.cmdId]; ok {
		command.Duration = time.Since(*pp.timer)

		if packetType(p) == responseErr {
			errorMsg, _ := decodeErrResponse(p)
			command.Error = errorMsg
			command.Result = responseErr
		} else {
			command.Error = ""
			command.Result = responseOk
		}

		pp.queryChan <- *command
		delete(pp.toRename, *pp.cmdId)
	}

	return len(p), nil
}

// NewProxyServer returns instance of ProxyServer
func NewProxyServer(mysqlAddr string, proxyAddr string) *ProxyServer {
	return &ProxyServer{
		mysqlAddr:   mysqlAddr,
		proxyAddr:   proxyAddr,
		Commands:    make(chan Command),
		Connections: make(chan Connection),
	}
}

// ProxyServer implements server for capturing and forwarding MySQL traffic.
type ProxyServer struct {
	mysqlAddr   string
	proxyAddr   string
	Commands    chan Command
	Connections chan Connection
}

// Run starts accepting TCP connection and forwarding it to MySQL server.
// Each incoming TCP connection is handled in own goroutine.
func (p *ProxyServer) Run() {
	listener, err := net.Listen("tcp", p.proxyAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer listener.Close()

	var connId uint

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Print(err.Error())
		}

		connId++
		go p.handleConnection(client, connId)
	}
}

// handleConnection ...
func (p *ProxyServer) handleConnection(client net.Conn, connId uint) {
	defer client.Close()

	// New connection to MySQL is made per each incoming TCP request to ProxyServer server.
	server, err := net.Dial("tcp", p.mysqlAddr)
	if err != nil {
		log.Print(err.Error())
		return
	}
	defer server.Close()

	defer func() { p.Connections <- Connection{connId, connStateFinished} }()

	var queryId uint
	var timer time.Time
	var commandsMap = make(map[uint]*Command)

	// Copy bytes from client to server and requestParser
	go io.Copy(io.MultiWriter(server, &RequestPacketParser{connId, &queryId, &timer, commandsMap}), client)

	// Copy bytes from server to client and responseParser
	io.Copy(io.MultiWriter(client, &ResponsePacketParser{connId, &queryId, p.Commands, &timer, commandsMap}), server)
}
