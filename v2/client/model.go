package client

import (
	"context"
	"crypto/tls"

	"github.com/lazybark/go-helpers/semver"
)

var ver = semver.Ver{
	Major:       2,
	Minor:       0,
	Patch:       2,
	Stable:      true,
	ReleaseNote: "not production tested",
}

type Client struct {
	//conn is the connection that will be used to read and write bytes
	conn *tls.Conn

	//host is the remote server to connect
	host string

	//isClosed is true when there is no connection or the connection was broken
	isClosed bool

	//isClosedWithError becomes true in case client has closed the connection due to an error
	isClosedWithError bool

	//conf points to client config
	conf *Config

	//ErrChan is the channel to send errors into external routine
	ErrChan chan error

	//ClientDoneChan is the channel to recieve client stopping command
	ClientDoneChan chan bool

	//MessageChan channel to notify external routine about new messages
	MessageChan chan *Message

	//bs holds total bytes sent by client
	bs int

	//br holds total bytes recieved by client in connection
	br int

	//errors holds total number of errors occured in connection
	errors int

	//connCount holds total number of successfull conections of the client
	connCount int

	//bytesLeft holds extra bytes that were read from stream after terminator occured, but end of buffer was not reached
	bytesLeft []byte

	//ctx is the connection context
	ctx    context.Context
	Cancel context.CancelFunc
}

type Config struct {
	//SuppressErrors stops client from sending errors into ErrChan.
	//Does not include fatal errors during startup.
	SuppressErrors bool

	//MaxMessageSize sets max length of one message in bytes.
	//If >0 and limit is reached, connection will be closed with an error.
	//
	//Note that if MaxMessageSize is > than reading buffer and MaxMessageSize reached,
	//it will not close connection until buffer is full or message terminator occurs.
	MaxMessageSize int

	//MessageTerminator sets byte value that marks message end in the stream.
	//Works for both incoming and outgoing messages
	MessageTerminator byte

	//BufferSize regulates buffer length to read incoming message. Default value is 128
	BufferSize int

	//DropOldStats = true will make client to set all sent/recieved bytes & errors to zero before opening new connection
	DropOldStats bool
}

// New creates new Client with specified config or default parameters
func New(conf *Config) *Client {
	c := new(Client)
	c.ErrChan = make(chan error)
	c.ClientDoneChan = make(chan bool)
	c.MessageChan = make(chan *Message)

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

// Version returns app version
func (c *Client) Version() semver.Ver { return ver }

// close closes connection and sets internal client vars to stop values
func (c *Client) close(err bool) error {
	if err {
		c.isClosedWithError = true
	}
	c.Cancel()
	c.isClosed = true
	return c.conn.Close()
}

// addRecBytes adds number to count of total recieved bytes
func (c *Client) addRecBytes(n int) { c.br += n }

// Stats returns cleint stats
func (c *Client) Stats() (int, int, int) { return c.bs, c.br, c.errors }

// Stats returns cleint stats
func (c *Client) ErrorsCount() int { return c.errors }

// Stats returns cleint stats
func (c *Client) RecievedBytes() int { return c.br }

// Stats returns cleint stats
func (c *Client) SentBytes() int { return c.bs }

// Close stops client and closes connection without error
func (c *Client) Close() error { return c.close(false) }
