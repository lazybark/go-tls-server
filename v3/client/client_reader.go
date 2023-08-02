package client

import (
	"fmt"

	"github.com/lazybark/go-tls-server/v3/conn"
)

// Reader infinitely reads messages from opened connection
func (c *Client) Reader() {
	for {
		if c.conn.Closed() {
			return
		}
		b, n, err := c.conn.ReadWithContext(c.conf.BufferSize, c.conf.MaxMessageSize, c.conf.MessageTerminator)
		if err != nil {
			if !c.conf.SuppressErrors {
				c.ErrChan <- fmt.Errorf("[Reader] error reading from %s -> %w", c.host, err)
			}
			return
		}

		//Check in case we read 0 bytes
		if n > 0 {
			c.MessageChan <- conn.NewMessage(c.conn, n, b)
		}
	}
}
