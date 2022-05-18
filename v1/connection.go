package v1

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

//Connection represents incoming client connection
type Connection struct {
	//id is an unique string to define connection
	id string
	//connectedAt time of connection init
	connectedAt time.Time
	//addr is the remote address of client
	addr net.Addr
	//conn is the connection interface that reads and writes bytes
	conn net.Conn
	//isClosed = true means that connection was closed and soon will be dropped from pool
	isClosed bool

	closedAt time.Time
	//lastAct updates every time there was any action in connection
	lastAct time.Time

	//ctx is the connection context
	ctx    context.Context
	cancel context.CancelFunc
	//bytesLeft holds extra bytes that were read from stream after terminator occured, but end of buffer was not reached
	bytesLeft []byte

	//bs holds total bytes sent by server in connection
	bs int
	//br holds total bytes sent by client in connection
	br int
	//errors holds total number of errors occured in connection
	errors int
}

//ConnectedAt returns time the connection was init
func (c *Connection) ConnectedAt() time.Time {
	return c.connectedAt
}

//Online returns duration of the connection
func (c *Connection) Online() time.Duration {
	if c.isClosed {
		return c.closedAt.Sub(c.connectedAt)
	}
	return time.Since(c.connectedAt)
}

//Address returns remote address of client
func (c *Connection) Address() net.Addr {
	return c.addr
}

//Closed returns true if the connection was closed
func (c *Connection) Closed() bool {
	return c.isClosed
}

//Id returns connection ID in pool
func (c *Connection) Id() string {
	return c.id
}

//close closes the connection with remote and sets isClosed as true
func (c *Connection) close() error {
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("[Connection][Close] error: %w", err)
	}
	c.isClosed = true
	return nil
}

//addRecBytes adds number to count of total recieved bytes
func (c *Connection) addRecBytes(n int) {
	c.br += n
}

//ReadContext reads bytes from connection until Terminator or error occurs or context is done.
//It can be used to read with timeout or any other way to break reader.
//
//Usual readers are vulnerable to routine-leaks, so this way is more confident.
func (c *Connection) readWithContext(buffer, maxSize int, terminator byte) ([]byte, int, error) {
	//Using c.conn.SetReadDeadline(time) in that case will make connection process less flexible.
	//Instead, checking ctx gives us a way to handle timeouts by the server itself.
	//We can, for example, close connection after some inactivity period by checking c.lastAct

	var rb []byte
	//Appending bytes that left from prev message in case terminator was not the last byte
	if len(c.bytesLeft) > 0 {
		rb = append(rb, c.bytesLeft...)
		c.bytesLeft = []byte{}
	}
	//Length of current read
	read := 0
	defer c.addRecBytes(read)
	//Read buffer with server-defined size
	b := make([]byte, buffer)
	for {
		select {
		case <-time.After(5 * time.Second):
			// Break by context
			return nil, read, fmt.Errorf("[ReadWithContext] reader closed by context")
		default:
			n, err := c.conn.Read(b)
			if err != nil {
				c.errors++
				return nil, read, fmt.Errorf("[ReadWithContext] reading error: %w", err)
			}
			read += n
			c.lastAct = time.Now()
			//We check every byte searching for terminator
			for num, by := range b[:n] {
				if by == terminator {
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
			if maxSize > 0 && read >= maxSize {
				c.errors++
				return nil, read, fmt.Errorf("[ReadWithContext] message size limits reached")
			}
			rb = append(rb, b[:n]...)

			if err == io.EOF {
				c.errors++
				c.close()
				return nil, read, fmt.Errorf("[ReadWithContext] stream closed")
			}
		}
	}
}

//SendByte sends bytes to remote by writing directrly into connection interface
func (c *Connection) SendByte(b []byte, term byte) (int, error) {
	b = append(b, term)
	bs, err := c.conn.Write(b)
	c.bs += bs
	c.lastAct = time.Now()
	if err != nil {
		c.errors++
		return bs, fmt.Errorf("[SendByte] error writing response: %w", err)
	}
	return bs, nil
}

//SendString converts s into byte slice and calls to SendByte
func (c *Connection) SendString(s string, term byte) (int, error) {
	return c.SendByte([]byte(s), term)
}
