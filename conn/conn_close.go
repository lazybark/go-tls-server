package conn

import (
	"fmt"

	"github.com/lazybark/go-helpers/npt"
)

// Closed returns true if the connection was closed.
func (c *Connection) Closed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.isClosed
}

// Close forsibly closes the connection. Active reader may still return bytes read between message start and connection close.
//
// IMPORTANT: 'close_notify' exchange is built on lower logic levels, but attempt to read/write with closed connection
// is still possible and will return error. If there is a risk that your app may do so, then you may need to use
// some flags to mark closed connections and avoid possible errors.
func (c *Connection) Close() error {
	if c.Closed() {
		return nil
	}

	return c.close()
}

// close marks connection as closed, but TLS will be closed by reader
func (c *Connection) close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cancel()
	c.isClosed = true
	c.closedAt = npt.Now()

	return nil
}

// closeTLS closes the TLS connection itself. To avoid data race it should be called by the reader function
func (c *Connection) closeTLS() error {
	err := c.tlsConn.Close()
	if err != nil {
		return fmt.Errorf("[Connection][closeTLS] %w", err)
	}

	return nil
}
