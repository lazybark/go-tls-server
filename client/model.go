package client

import (
	"fmt"
	"sync"

	"github.com/lazybark/go-helpers/semver"
	"github.com/lazybark/go-tls-server/conn"
)

var ver = semver.Ver{
	Major:       3,
	Minor:       1,
	Patch:       0,
	Stable:      false,
	ReleaseNote: "beta",
}

// Client is the TLS client managing one single connection & statistics.
//
// IMPORTANT: it has all stats & methods assigned to the Client struct itself, not to separate Connection struct
// as server does.
type Client struct {
	// conn is the connection that will be used to read and write bytes.
	conn *conn.Connection

	// host is the remote server to connect.
	host string

	// isClosed is true when there is no connection or the connection was broken.
	isClosed bool

	// isClosedWithError becomes true in case client has closed the connection due to an error.
	isClosedWithError bool

	// conf points to client config.
	conf *Config

	// errChan is the channel to send errors into external routine.
	errChan chan error

	// ClientDoneChan is the channel to receive client stopping command from external routine.
	// It's not used by default and exist to provide flexibility for bigger apps that will use client.
	ClientDoneChan chan bool

	// messageChan channel to notify external routine about new messages.
	messageChan chan *conn.Message

	// connCount holds total number of successful conections of the client.
	connCount int

	// mu is used to set client closed to state. It's not protecting stat variables which means
	// reading and/or writing should not be done concurrently.
	mu *sync.RWMutex
}

// ErrChan returns clients's error channel to read only.
func (c *Client) ErrChan() <-chan error {
	return c.errChan
}

// MessageChan returns clients's message channel to read only.
func (c *Client) MessageChan() <-chan *conn.Message {
	return c.messageChan
}

// Stats returns number of bytes sent/receive + number of errors.
func (c *Client) Stats() (sent, received, errors int) { return c.conn.Stats() }

// Version returns app version.
func (c *Client) Version() semver.Ver { return ver }

// close closes connection and sets internal client vars to stop values.
func (c *Client) close(withError bool) error {
	c.mu.Lock()
	c.isClosedWithError = withError
	c.isClosed = true
	c.mu.Unlock()

	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("[close] %w", err)
	}

	return nil
}

// Close stops client and closes connection without error.
//
// Important: it does not close the message & error channels as there still might be some leftovers got by client.
// You may want to close them manually if your app design needs it.
func (c *Client) Close() error { return c.close(false) }

// Close stops client and closes connection WITH error.
//
// Important: it does not close the message & error channels as there still might be some leftoverss got by client.
// You may want to close them manually if your app design needs it.
func (c *Client) CloseWithError() error { return c.close(true) }

// Closed returns true if client was closed.
func (c *Client) Closed() (cl bool) {
	c.mu.RLock()
	cl = c.isClosed
	c.mu.RUnlock()

	return cl
}

// ClosedWithError returns true if client was closed with error.
func (c *Client) ClosedWithError() (cl bool) {
	c.mu.RLock()
	cl = c.isClosedWithError
	c.mu.RUnlock()

	return
}
