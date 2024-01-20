package client

import "fmt"

// SendByte sends bytes to remote by writing directrly into connection interface.
func (c *Client) SendByte(b []byte) (int, error) {
	count, err := c.conn.SendByte(b)
	if err != nil {
		return count, c.FormatError(fmt.Errorf("[SendByte]: %w", err))
	}

	return count, nil
}

// SendString converts s into byte slice and calls to SendByte.
func (c *Client) SendString(s string) (int, error) {
	count, err := c.conn.SendString(s)
	if err != nil {
		return count, c.FormatError(fmt.Errorf("[SendString]: %w", err))
	}

	return count, nil
}
