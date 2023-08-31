package client

import (
	"sync"

	"github.com/lazybark/go-tls-server/conn"
)

// New creates new Client with specified config or default parameters
func New(conf *Config) *Client {
	c := new(Client)
	c.errChan = make(chan error, 3)
	c.ClientDoneChan = make(chan bool)
	c.messageChan = make(chan *conn.Message, 10)
	c.mu = &sync.RWMutex{}

	if conf == nil {
		conf = &Config{}
		//Dropping all stats is the default behaviour
		conf.DropOldStats = true
	}
	//Default terminator is the newline
	if conf.MessageTerminator == 0 {
		conf.MessageTerminator = '\n'
	}
	//Default buffer is 128 B
	if conf.BufferSize == 0 {
		conf.BufferSize = 128
	}
	c.conf = conf

	return c
}
