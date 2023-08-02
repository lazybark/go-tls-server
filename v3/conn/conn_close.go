package conn

import (
	"fmt"

	"github.com/lazybark/go-helpers/npt"
)

// Closed returns true if the connection was closed
func (c *Connection) Closed() (closed bool) {
	c.mu.RLock()
	closed = c.isClosed
	c.mu.RUnlock()

	return
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

// close closes the connection with remote
func (c *Connection) close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cancel()
	err := c.tlsConn.Close()
	if err != nil {
		return fmt.Errorf("[Connection][close] %w", err)
	}
	c.isClosed = true
	c.closedAt = npt.Now()

	return nil
}
