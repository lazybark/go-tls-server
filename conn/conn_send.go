package conn

import "fmt"

// SendByte sends bytes to remote by writing directrly into connection interface.
func (c *Connection) SendByte(bytesToSend []byte) (int, error) {
	bytesToSend = append(bytesToSend, c.messageTerminator)
	sentCount, err := c.tlsConn.Write(bytesToSend)

	c.AddSentBytes(sentCount)
	c.SetLastAct()

	if err != nil {
		c.AddErrors(1)

		return sentCount, fmt.Errorf("[SendByte] error writing response: %w", err)
	}

	return sentCount, nil
}

// SendString converts s into byte slice and calls to SendByte.
func (c *Connection) SendString(s string) (int, error) { return c.SendByte([]byte(s)) }
