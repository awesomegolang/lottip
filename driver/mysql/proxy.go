package mysql

import (
	"io"
	"log"
	"net"
	"time"
)

type requestWriter struct {
	connId      uint
	cmdId       *uint
	commandsMap map[uint]*Command
}

func (pp *requestWriter) Write(p []byte) (n int, err error) {
	*pp.cmdId++

	switch packetType(p) {
	case comStmtPrepare:
	case comQuery:
		decoded, err := decodeQueryRequest(p)
		if err == nil {
			pp.commandsMap[*pp.cmdId] = &Command{
				ConnId:     pp.connId,
				CmdId:      *pp.cmdId,
				Query:      decoded.Query,
				Executable: false,
				StartedAt: time.Now(),
			}
		}
	}

	return len(p), nil
}

type responseWriter struct {
	cmdId       *uint
	queryChan   chan Command
	commandsMap map[uint]*Command
}

func (pp *responseWriter) Write(p []byte) (n int, err error) {
	if command, ok := pp.commandsMap[*pp.cmdId]; ok {
		command.Duration = time.Since(command.StartedAt)

		if packetType(p) == responseErr {
			errorMsg, _ := decodeErrResponse(p)
			command.Error = errorMsg
			command.Result = responseErr
		} else {
			command.Error = ""
			command.Result = responseOk
		}

		pp.queryChan <- *command
		delete(pp.commandsMap, *pp.cmdId)
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

	p.Connections <- Connection{connId, connStateStarted}
	defer func() { p.Connections <- Connection{connId, connStateFinished} }()

	var cmdId uint
	var commandsMap = make(map[uint]*Command)

	// Copy bytes from client to server and requestParser
	go io.Copy(io.MultiWriter(server, &requestWriter{connId, &cmdId, commandsMap}), client)

	// Copy bytes from server to client and responseParser
	io.Copy(io.MultiWriter(client, &responseWriter{&cmdId, p.Commands, commandsMap}), server)
}
