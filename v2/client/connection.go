package client

import (
	"fmt"
	"io"
)

// ReadWithContext reads bytes from connection until Terminator or error occurs or context is done.
// It can be used to read with timeout or any other way to break reader.
//
// Usual readers are vulnerable to routine-leaks, so this way is more confident.
func (c *Client) ReadWithContext() ([]byte, int, error) {
	//Using c.conn.SetReadDeadline(time) in that case will make connection process less flexible.
	//Instead, checking ctx gives us a way to handle timeouts by the client itself.
	//We can, for example, close connection after some inactivity period by checking c.lastAct

	var rb []byte
	//Appending bytes that left from prev message in case terminator was not the last byte
	if len(c.bytesLeft) > 0 {
		rb = append(rb, c.bytesLeft...)
		c.bytesLeft = []byte{}
	}
	//Length of current read
	read := 0
	//Read buffer with defined size
	b := make([]byte, c.conf.BufferSize)
	for {
		select {
		case <-c.ctx.Done():
			// Break by context
			return nil, read, fmt.Errorf("[ReadWithContext] reader closed by context")
		default:
			n, err := c.conn.Read(b)
			if err != nil {
				c.errors++
				if err == io.EOF {
					return nil, read, fmt.Errorf("[ReadWithContext] stream closed")
				}
				if c.ctx.Done() != nil {
					return nil, read, fmt.Errorf("[ReadWithContext] reader closed by context")
				}
				return nil, read, fmt.Errorf("[ReadWithContext] reading error: %w", err)
			}
			read += n
			c.addRecBytes(n)
			//We check every byte searching for terminator
			for num, by := range b[:n] {
				if by == c.conf.MessageTerminator {
					rb = append(rb, b[:num]...)
					//We collect extra bytes in case there is something left from prev message and pass on to next one
					//This can happen in cases when client sends data in a stream-way, not portionally
					//These bytes will be picked up with next trigger of reader
					if len(b[num:n]) > 0 {
						c.bytesLeft = b[num:n]
					}
					return rb, read, nil
				}
			}
			if c.conf.MaxMessageSize > 0 && read > c.conf.MaxMessageSize {
				c.errors++
				return nil, read, fmt.Errorf("[ReadWithContext] message size limits reached")
			}
			rb = append(rb, b[:n]...)
		}
	}
}

// SendByte sends bytes to remote by writing directrly into connection
func (c *Client) SendByte(b []byte) error {
	b = append(b, c.conf.MessageTerminator)
	bs, err := c.conn.Write(b)
	if err != nil {
		return fmt.Errorf("[SendByte] error writing response: %w", err)
	}
	c.bs += bs
	return nil
}

// SendString converts s into byte slice and calls to SendByte
func (c *Client) SendString(s string) error { return c.SendByte([]byte(s)) }
