package v1

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"
)

var ver = "1.0.4"

type Server struct {
	timeStart time.Time
	//host = hostname of the server
	host string
	//connPool is a map of pointers to connections.
	//
	//In this case pointers are used to increase code readability and number of ops
	//needed to change conn state
	connPool map[string]*Connection
	//connPoolMutex controls connPool
	connPoolMutex sync.RWMutex

	//listener is the interface that listens for new connections
	listener net.Listener

	//tlsConfig points to tls listener config
	tlsConfig *tls.Config
	//sConfig points to server config
	sConfig *Config

	//ErrChan is the channel to send errors into external routine
	ErrChan chan error
	//ServerDoneChan is the channel to recieve server stopping command
	ServerDoneChan chan bool
	//ConnChan is the channel to notify external routine about new connection
	ConnChan chan *Connection
	//MessageChan channel to notify external routine about new messages
	MessageChan chan *Message

	//Stat keeps connections stat by date
	Stat map[string]Stat
}

//New initializes server instance and makes it completely ready to listen for connections
func New(host string, cert string, key string, conf *Config) (*Server, error) {
	s := new(Server)
	s.timeStart = time.Now()
	s.host = host
	s.ErrChan = make(chan error)
	s.ServerDoneChan = make(chan bool)
	s.ConnChan = make(chan *Connection)
	s.MessageChan = make(chan *Message)
	s.connPool = make(map[string]*Connection)
	s.Stat = make(map[string]Stat)
	s.connPoolMutex = sync.RWMutex{}

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
	//Start server admin
	go s.adminRoutine()

	return s, nil
}

//StartedAt returns starting time
func (s *Server) StartedAt() time.Time { return s.timeStart }

//Online returns time online
func (s *Server) Online() time.Duration { return time.Since(s.timeStart) }

//ActiveConnetions returns number of active connections
func (s *Server) TotalConnetions() int {
	a := 0
	for _, c := range s.connPool {
		if !c.isClosed {
			a++
		}
	}
	return a
}

//Version returns app version
func (s *Server) Version() string { return ver }

//CloseConnection is the only correct way to close connection.
//It changes conn state in pool and then calls to c.close
func (s *Server) CloseConnection(c *Connection) error {
	c.isClosed = true
	return c.close()
}

//addToPool adds connection to fool for controlling
func (s *Server) addToPool(c *Connection) {
	s.connPoolMutex.Lock()
	s.connPool[c.id] = c
	s.connPoolMutex.Unlock()
}

//remFromPool removes connection pointer from pool, so it becomes unavailable to reach
func (s *Server) remFromPool(c *Connection) {
	s.connPoolMutex.Lock()
	delete(s.connPool, c.id)
	s.connPoolMutex.Unlock()
}

//adminRoutine controls server behaviour: drops closed connections, closes inactive ones and stops the server in case s.ServerDoneChan
func (s *Server) adminRoutine() {
	for {
		select {
		//Once per hour clean up old & close inactive connections
		case <-time.After(time.Hour):
			for _, c := range s.connPool {
				//If conn is closed and time now is already after the moment it should be deleted permanently
				if c.isClosed && !time.Now().Before(c.ClosedAt().Add(time.Minute*time.Duration(s.sConfig.KeepOldConnections))) {
					s.remFromPool(c)
					continue
				}
				//If it's not closed, but it's been a 'KeepInactiveConnections' time after
				if !c.isClosed && s.sConfig.KeepInactiveConnections > 0 && !time.Now().Before(c.LastAct().Add(time.Minute*time.Duration(s.sConfig.KeepInactiveConnections))) {
					s.CloseConnection(c)
				}
			}
		//In case server needs to be stopped - close all connections
		case d := <-s.ServerDoneChan:
			if d {
				err := s.listener.Close()
				if err != nil && !s.sConfig.SuppressErrors {
					s.ErrChan <- fmt.Errorf("[Server][Listen] error closing listener: %w", err)
				}
				for _, c := range s.connPool {
					err := s.CloseConnection(c)
					if err != nil && !s.sConfig.SuppressErrors {
						s.ErrChan <- fmt.Errorf("[Server][adminRoutine] error closing connection %s -> %w", c.id, err)
					}
				}
			}
		}
	}
}

//SendByte calls to c.SendByte and adds sent bytes to Stat
func (s *Server) SendByte(c *Connection, b []byte) error {
	n, err := c.SendByte(b, s.sConfig.MessageTerminator)
	s.addSentBytes(n)
	if err != nil {
		s.addErrors(1)
	}
	return err
}

//SendString calls to c.SendString and adds sent bytes to Stat
func (s *Server) SendString(c *Connection, str string) error {
	n, err := c.SendString(str, s.sConfig.MessageTerminator)
	s.addSentBytes(n)
	if err != nil {
		s.addErrors(1)
	}
	return err
}
