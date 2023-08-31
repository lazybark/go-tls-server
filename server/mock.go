package server

import (
	"sync"
	"time"

	"github.com/lazybark/go-tls-server/conn"
)

func GetEmptyTestServer() *Server {
	s := new(Server)
	s.timeStart = time.Now()
	s.host = "localhost"
	s.errChan = make(chan error)
	s.serverDoneChan = make(chan bool)
	s.connChan = make(chan *conn.Connection)
	s.connPool = make(map[string]*conn.Connection)
	s.stat = make(map[string]Stat)
	s.statOverall = &Stat{}
	s.connPoolMutex = sync.RWMutex{}
	s.ver = ver

	return s
}
