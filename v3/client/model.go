package client

import (
	"github.com/lazybark/go-helpers/semver"
	"github.com/lazybark/go-tls-server/v3/conn"
)

var ver = semver.Ver{
	Major:       3,
	Minor:       0,
	Patch:       2,
	Stable:      false,
	ReleaseNote: "beta",
}

// Client is the TLS client managing one single connection & statistics.
//
// IMPORTANT: it has all stats & methods assigned to the Client struct itself, not to separate Connection struct
// as server does.
type Client struct {
	//conn is the connection that will be used to read and write bytes
	conn *conn.Connection

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

	//ClientDoneChan is the channel to recieve client stopping command from external routine.
	//It's not used by default and exist to provide flexibility for bigger apps that will use client
	ClientDoneChan chan bool

	//MessageChan channel to notify external routine about new messages
	MessageChan chan *conn.Message

	//connCount holds total number of successfull conections of the client
	connCount int
}

// Stats returns number of bytes sent/receive + number of errors
func (c *Client) Stats() (sent, received, errors int) { return c.conn.Stats() }

// Version returns app version
func (c *Client) Version() semver.Ver { return ver }

// close closes connection and sets internal client vars to stop values
func (c *Client) close(withError bool) error {
	c.isClosedWithError = withError
	c.isClosed = true
	return c.conn.Close()
}

// Close stops client and closes connection without error
func (c *Client) Close() error { return c.close(false) }

// Close stops client and closes connection WITH error
func (c *Client) CloseWithError() error { return c.close(true) }

// Closed returns true if client was closed
func (c *Client) Closed() bool { return c.isClosed }

// ClosedWithError returns true if client was closed with error
func (c *Client) ClosedWithError() bool { return c.isClosedWithError }
