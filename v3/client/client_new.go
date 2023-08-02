package client

import (
	"context"

	"github.com/lazybark/go-tls-server/v3/conn"
)

// New creates new Client with specified config or default parameters
func New(conf *Config) *Client {
	c := new(Client)
	c.ErrChan = make(chan error)
	c.ClientDoneChan = make(chan bool)
	c.MessageChan = make(chan *conn.Message)

	ctx, cancel := context.WithCancel(context.Background())
	c.Cancel = cancel
	c.ctx = ctx

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
