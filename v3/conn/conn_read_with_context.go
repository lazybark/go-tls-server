package conn

import (
	"fmt"
	"io"
)

// readWithContext reads bytes from connection until Terminator / error occurs or context is done.
// It can be used to read with timeout or any other way to break reader.
// Usual readers are vulnerable to routine-leaks, so this way is more confident.
//
// IMPORTANT: if EOF or context deadline appear, readWithContext will mark connection as 'closed'.
// Other errors should be treated manually by external code.
// In all cases method will return last bytes read
func (c *Connection) ReadWithContext(buffer, maxSize int, terminator byte) ([]byte, int, error) {
	if c.Closed() {
		return nil, 0, fmt.Errorf("[ReadWithContext] %w", ErrReaderAlreadyClosed)
	}

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
	defer func(read *int) { c.addRecBytes(*read) }(&read)
	//Read buffer with server-defined size
	b := make([]byte, buffer)
	for {
		select {
		case <-c.ctx.Done():
			// Break by context
			_ = c.closeTLS() //We close TLS only by reader

			return nil, read, nil
		default:
			n, err := c.tlsConn.Read(b)
			if err != nil {
				if err == io.EOF {
					return nil, read, fmt.Errorf("[ReadWithContext] %w", ErrStreamClosed)
				}
				if c.ctx.Done() != nil {
					_ = c.closeTLS() //We close TLS only by reader

					return nil, read, nil
				}
				c.addErrors(1)
				//The connecton is not closed yet in this case!
				//Client code should decide if they want to close or try to read next bytes
				return nil, read, fmt.Errorf("[ReadWithContext] reading error: %w", err)
			}
			read += n

			c.setLastAct()

			if maxSize > 0 && read > maxSize {
				c.addErrors(1)
				return nil, read, fmt.Errorf("[ReadWithContext] %w (read %v of max %v)", ErrMessageSizeLimit, read, maxSize)
			}
			//We check every byte searching for terminator
			for num, by := range b[:n] {
				if by == terminator {
					rb = append(rb, b[:num]...)
					//We collect extra bytes in case there is something left from prev message and pass on to next one
					//This can happen in cases when client sends data in a stream-way, not portionally
					//These bytes will be picked up with next trigger of reader as if they were sent with next message itself
					if len(b[num:n]) > 0 {
						c.bytesLeft = b[num:n]
					}
					return rb, read, nil
				}
			}

			rb = append(rb, b[:n]...)
		}
	}
}
