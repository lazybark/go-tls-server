package conn

import "time"

// Stats returns Connection stats
func (c *Connection) Stats() (sent, received, errors int) { return c.bs, c.br, c.errors }

// DropOldStats sets bytes recieved, sent and error count to zero
func (c *Connection) DropOldStats() {
	c.mu.Lock()
	c.br = 0
	c.bs = 0
	c.errors = 0
	c.mu.Unlock()
}

// ConnectedAt returns time the connection was init
func (c *Connection) ConnectedAt() time.Time { return c.connectedAt.Time() }

// ConnectedAt returns time the connection was init
func (c *Connection) ClosedAt() time.Time { return c.closedAt.Time() }

// ConnectedAt returns time the connection was init
func (c *Connection) LastAct() time.Time { return c.lastAct.Time() }

// Online returns duration of the connection
func (c *Connection) Online() time.Duration {
	if c.isClosed {
		return c.ClosedAt().Sub(c.ConnectedAt())
	}
	return time.Since(c.ConnectedAt())
}

// addRecBytes adds number to count of total recieved bytes
func (c *Connection) addRecBytes(n *int) { c.br += *n }
