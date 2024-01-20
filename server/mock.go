package server

import (
	"sync"
	"time"

	"github.com/lazybark/go-tls-server/conn"
)

func GetEmptyTestServer() *Server {
	server := new(Server)
	server.timeStart = time.Now()
	server.host = "localhost"
	server.errChan = make(chan error)
	server.serverDoneChan = make(chan bool)
	server.connChan = make(chan *conn.Connection)
	server.connPool = make(map[string]*conn.Connection)
	server.stat = make(map[string]Stat)
	server.statOverall = new(Stat)
	server.connPoolMutex = sync.RWMutex{}

	return server
}
