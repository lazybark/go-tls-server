package server

import (
	"crypto/tls"
	"fmt"
	"log"
)

// Listen runs listener interface implementations and accepts connections
func (s *Server) Listen(port string) {
	l, err := tls.Listen("tcp", ":"+port, s.tlsConfig)
	if err != nil {
		log.Fatal(fmt.Errorf("[Server][Listen] error listening: %w", err))
	}

	s.listener = l

	defer func() {
		err := l.Close()
		if err != nil && !s.sConfig.SuppressErrors {
			s.ErrChan <- fmt.Errorf("[Server][Listen] error closing listener: %w", err)
		}
	}()

	for {
		//Accept the connection
		conn, err := l.Accept()
		if err != nil && !s.sConfig.SuppressErrors {
			s.ErrChan <- fmt.Errorf("[Server][Listen] error accepting connection from %v: %w", conn.RemoteAddr(), err)
		}

		c, err := NewConnection(conn.RemoteAddr(), conn, s.sConfig.MessageTerminator)
		if err != nil && !s.sConfig.SuppressErrors {
			s.ErrChan <- fmt.Errorf("[Server][Listen] error making connection Id for %v: %w", conn.RemoteAddr(), err)
		}

		//Add to pool
		s.addToPool(c)
		//Notify outer routine
		if s.sConfig.NotifyAboutNewConnections {
			s.ConnChan <- c
		}
		//Wait for new messages
		go s.recieve(c)
	}
}

// recieve endlessy reads incoming stream and delivers messages to recievers outside server routine.
// It uses ReadWithContext, so execution can be manually stopped by calling c.cancel on specific connection.
// In that case (or if any error occurs) method will trigger s.CloseConnection to break connection too
func (s *Server) recieve(c *Connection) {
	for {
		b, n, err := c.readWithContext(s.sConfig.BufferSize, s.sConfig.MaxMessageSize, s.sConfig.MessageTerminator)
		if err != nil {
			if !s.sConfig.SuppressErrors {
				s.ErrChan <- fmt.Errorf("[Server][recieve] error reading from %s: %w", c.id, err)
			}
			s.CloseConnection(c)
			return
		}
		s.addRecBytes(n)
		c.MessageChan <- &Message{conn: c, length: n, bytes: b}
	}
}
