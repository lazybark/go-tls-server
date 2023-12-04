package client

// SendByte sends bytes to remote by writing directrly into connection interface.
func (c *Client) SendByte(b []byte) (int, error) {
	return c.conn.SendByte(b)
}

// SendString converts s into byte slice and calls to SendByte.
func (c *Client) SendString(s string) (int, error) {
	return c.conn.SendString(s)
}
