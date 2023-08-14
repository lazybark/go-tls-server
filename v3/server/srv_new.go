package server

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/lazybark/go-tls-server/v3/conn"
)

// New initializes server instance and makes it completely ready to listen for connections
func New(host string, cert string, key string, conf *Config) (*Server, error) {
	s := new(Server)
	s.timeStart = time.Now()
	s.host = host
	s.ErrChan = make(chan error)
	s.ServerDoneChan = make(chan bool)
	s.ConnChan = make(chan *conn.Connection)
	s.connPool = make(map[string]*conn.Connection)
	s.stat = make(map[string]Stat)
	s.statOverall = new(Stat)
	s.connPoolMutex = sync.RWMutex{}
	s.ver = ver

	if conf == nil {
		conf = &Config{}
	}
	//Default terminator is the newline
	if conf.MessageTerminator == 0 {
		conf.MessageTerminator = '\n'
	}

	//Default buffer is 128 B
	if conf.BufferSize == 0 {
		conf.BufferSize = 128
	}
	//KeepOldConnections by default is 24 hours
	if conf.KeepOldConnections == 0 {
		conf.KeepOldConnections = 1440
	}
	//KeepInactiveConnections by default is 72 hours
	if conf.KeepInactiveConnections == 0 {
		conf.KeepInactiveConnections = 4320
	}
	s.sConfig = conf

	var tlsConfig *tls.Config
	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("[Server] error getting key pair: %w", err)
	}
	tlsConfig = &tls.Config{Certificates: []tls.Certificate{certificate}}
	s.tlsConfig = tlsConfig

	if conf.HttpStatMode {
		if conf.HttpStatAddr == "" {
			conf.HttpStatAddr = "localhost:3939"
		}
		s.resolver = chi.NewRouter()
	}

	//Start server admin
	go s.adminRoutine()

	return s, nil
}
