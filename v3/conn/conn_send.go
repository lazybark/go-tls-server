package conn

import "fmt"

// SendByte sends bytes to remote by writing directrly into connection interface
func (c *Connection) SendByte(b []byte) (int, error) {
	b = append(b, c.messageTerminator)
	bs, err := c.tlsConn.Write(b)

	c.addSentBytes(bs)
	c.setLastAct()

	if err != nil {
		c.addErrors(1)

		return bs, fmt.Errorf("[SendByte] error writing response: %w", err)
	}

	return bs, nil
}

// SendString converts s into byte slice and calls to SendByte
func (c *Connection) SendString(s string) (int, error) { return c.SendByte([]byte(s)) }
