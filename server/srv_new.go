package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/lazybark/go-tls-server/conn"
)

// New initializes server instance and makes it completely ready to listen for connections.
func New(host string, cert string, key string, conf *Config) (*Server, error) {
	s := new(Server)
	s.timeStart = time.Now()
	s.host = host
	s.errChan = make(chan error)
	s.serverDoneChan = make(chan bool)
	s.connChan = make(chan *conn.Connection)
	s.connPool = make(map[string]*conn.Connection)
	s.stat = make(map[string]Stat)
	s.statOverall = new(Stat)
	s.connPoolMutex = sync.RWMutex{}
	s.ver = ver

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.ctx = ctx

	if conf == nil {
		conf = &Config{}
	}
	// Default terminator is the newline.
	if conf.MessageTerminator == 0 {
		conf.MessageTerminator = '\n'
	}

	// Default buffer is 128 B.
	if conf.BufferSize == 0 {
		conf.BufferSize = 128
	}
	// KeepOldConnections by default is 24 hours.
	if conf.KeepOldConnections == 0 {
		conf.KeepOldConnections = 1440
	}
	// KeepInactiveConnections by default is 72 hours.
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

	// Start server admin.
	go s.adminRoutine()

	return s, nil
}
