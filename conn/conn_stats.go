package conn

import "time"

// Stats returns Connection stats
func (c *Connection) Stats() (sent, received, errors int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.bs, c.br, c.errors
}

// DropOldStats sets bytes received, sent and error count to zero
func (c *Connection) DropOldStats() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.br = 0
	c.bs = 0
	c.errors = 0
}

// ConnectedAt returns time the connection was init
func (c *Connection) ConnectedAt() time.Time { return c.connectedAt.Time() }

// ConnectedAt returns time the connection was init
func (c *Connection) ClosedAt() time.Time { return c.closedAt.Time() }

func (c *Connection) SetLastAct() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastAct.ToNow()
}

// ConnectedAt returns time the connection was init
func (c *Connection) LastAct() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.lastAct.Time()
}

// Online returns duration of the connection
func (c *Connection) Online() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		return c.ClosedAt().Sub(c.ConnectedAt())
	}

	return time.Since(c.ConnectedAt())
}

// AddRecBytes adds number to count of total received bytes
func (c *Connection) AddRecBytes(n int) {
	if n < 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.br += n
}

// AddSentBytes adds number to count of total sent bytes
func (c *Connection) AddSentBytes(n int) {
	if n < 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.bs += n
}

// AddErrors adds number to count of total errors
func (c *Connection) AddErrors(n int) {
	if n < 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.errors += n
}
