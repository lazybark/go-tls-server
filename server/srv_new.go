package server

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/lazybark/go-helpers/semver"
	"github.com/lazybark/go-tls-server/conn"
)

// New initializes server instance and makes it completely ready to listen for connections.
func New(host string, cert string, key string, conf *Config) (*Server, error) { //nolint: funlen // false alarm
	server := new(Server)
	server.timeStart = time.Now()
	server.host = host
	server.errChan = make(chan error)
	server.serverDoneChan = make(chan bool)
	server.connChan = make(chan *conn.Connection)
	server.connPool = make(map[string]*conn.Connection)
	server.stat = make(map[string]Stat)
	server.statOverall = new(Stat)
	server.connPoolMutex = sync.RWMutex{}
	server.mu = new(sync.Mutex)
	server.ver = semver.Ver{ //nolint:exhaustruct // false alarm
		Major:       3, //nolint:gomnd // false alarm
		Minor:       2, //nolint:gomnd // false alarm
		Patch:       0,
		Stable:      false,
		ReleaseNote: "beta",
	}

	if conf == nil {
		conf = new(Config)
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

	if conf.ErrorPrefix == "" {
		conf.ErrorPrefix = "TLS_SERVER"
	}

	server.sConfig = conf

	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, server.FormatError(fmt.Errorf("error getting key pair: %w", err))
	}

	server.tlsConfig = new(tls.Config)
	server.tlsConfig.Certificates = []tls.Certificate{certificate}
	server.tlsConfig.MinVersion = tls.VersionTLS12

	// Start server admin.
	go server.adminRoutine()

	return server, nil
}
