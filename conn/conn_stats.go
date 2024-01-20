package conn

import "time"

// Stats returns Connection stats.
func (c *Connection) Stats() (int, int, int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.bs, c.br, c.errors
}

// DropOldStats sets bytes received, sent and error count to zero.
func (c *Connection) DropOldStats() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.br = 0
	c.bs = 0
	c.errors = 0
}

// ConnectedAt returns time the connection was init.
func (c *Connection) ConnectedAt() time.Time { return c.connectedAt.Time() }

// ConnectedAt returns time the connection was init.
func (c *Connection) ClosedAt() time.Time { return c.closedAt.Time() }

func (c *Connection) SetLastAct() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastAct.ToNow()
}

// ConnectedAt returns time the connection was init.
func (c *Connection) LastAct() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.lastAct.Time()
}

// Online returns duration of the connection.
func (c *Connection) Online() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		return c.ClosedAt().Sub(c.ConnectedAt())
	}

	return time.Since(c.ConnectedAt())
}

// AddRecBytes adds number to count of total received bytes.
func (c *Connection) AddRecBytes(count int) {
	if count < 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.br += count
}

// AddSentBytes adds number to count of total sent bytes.
func (c *Connection) AddSentBytes(count int) {
	if count < 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.bs += count
}

// AddErrors adds number to count of total errors.
func (c *Connection) AddErrors(count int) {
	if count < 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.errors += count
}

// Sent returns total count of bytes sent into connection.
func (c *Connection) Sent() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.bs
}

// Received returns total count of bytes received from connection.
func (c *Connection) Received() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.br
}

// Sent returns total count of bytes sent into connection.
func (c *Connection) Errors() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.errors
}
