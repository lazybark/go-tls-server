package server

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/lazybark/go-tls-server/v3/conn"
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
		tlsConn, err := l.Accept()
		if err != nil && !s.sConfig.SuppressErrors {
			s.ErrChan <- fmt.Errorf("[Server][Listen] error accepting connection from %v: %w", tlsConn.RemoteAddr(), err)
		}

		c, err := conn.NewConnection(tlsConn.RemoteAddr(), tlsConn, s.sConfig.MessageTerminator)
		if err != nil && !s.sConfig.SuppressErrors {
			s.ErrChan <- fmt.Errorf("[Server][Listen] error making connection for %v: %w", tlsConn.RemoteAddr(), err)
		}

		//Add to pool
		s.addToPool(c)
		//Notify outer routine
		s.ConnChan <- c
		//Wait for new messages
		go s.recieve(c)
	}
}

// recieve endlessy reads incoming stream and delivers messages to recievers outside server routine.
// It uses ReadWithContext, so execution can be manually stopped by calling c.cancel on specific connection.
// In that case (or if any error occurs) method will trigger s.CloseConnection to break connection too
func (s *Server) recieve(c *conn.Connection) {
	for {
		if c.Closed() {
			return
		}
		b, n, err := c.ReadWithContext(s.sConfig.BufferSize, s.sConfig.MaxMessageSize, s.sConfig.MessageTerminator)
		if err != nil {
			if !s.sConfig.SuppressErrors {
				s.ErrChan <- fmt.Errorf("[Server][recieve] error reading from %s: %w", c.Id(), err)
			}
			s.CloseConnection(c)
			return
		}

		//Check in case we read 0 bytes
		if n > 0 {
			s.addRecBytes(n)
			c.MessageChan <- conn.NewMessage(c, n, b)
		}
	}
}
