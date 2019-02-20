package redis

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Command struct {
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at"`
}

type Monitor struct {
	host     string
	Commands chan Command
}

func NewMonitor(host string) *Monitor {
	return &Monitor{
		host:     host,
		Commands: make(chan Command),
	}
}

func (m *Monitor) Run() {
	conn, err := net.Dial("tcp", m.host)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	fmt.Fprint(conn, "MONITOR\n")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		m.Commands <- Command{scanner.Text(), time.Now()}
	}
}
