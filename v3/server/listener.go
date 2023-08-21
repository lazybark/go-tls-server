package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/lazybark/go-tls-server/v3/conn"
)

// Listen runs listener interface implementations and accepts connections
func (s *Server) Listen(port string) {
	l, err := tls.Listen("tcp", ":"+port, s.tlsConfig)
	if err != nil {
		log.Fatal(fmt.Errorf("[Server][Listen] error listening: %w", err))
	}

	s.listener = l

	//Start HTTP server
	if s.sConfig.HttpStatMode {
		s.setHTTPRoutes()

		go s.serveHTTP()
	}

	for {

		select {
		case <-s.ctx.Done():

			return

		default:
			//Accept the connection
			tlsConn, err := l.Accept()

			//The problem is that a listener can be closed during the listening. Then we get net.ErrClosed.
			//In this case we always ignore it, because doesn't matter why it's closed - this function is not for err processing
			//Error was handled somewhere else already. Or server was simply terminated
			//TO DO: think if there is more graceful way than errors.Is
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					continue
				}
				if !s.sConfig.SuppressErrors {
					s.ErrChan <- fmt.Errorf("[Server][Listen] error accepting connection: %w", err)
				}
			}

			//Just a precaution to avoid nil pointer dereference
			if tlsConn == nil {
				continue
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
			go s.receive(c)
		}
	}
}

// receive endlessly reads incoming stream and delivers messages to receivers outside server routine.
// It uses ReadWithContext, so execution can be manually stopped by calling c.cancel on specific connection.
// In that case (or if any error occurs) method will trigger s.CloseConnection to break connection too
func (s *Server) receive(c *conn.Connection) {
	for {
		if c.Closed() {
			return
		}
		b, n, err := c.ReadWithContext(s.sConfig.BufferSize, s.sConfig.MaxMessageSize, s.sConfig.MessageTerminator)
		if err != nil {
			if !s.sConfig.SuppressErrors {
				s.ErrChan <- fmt.Errorf("[Server][receive] error reading from %s: %w", c.Id(), err)
			}
			err := s.CloseConnection(c)
			if err != nil && !s.sConfig.SuppressErrors {
				s.ErrChan <- fmt.Errorf("[Server][receive] error closing connection: %w", err)
			}

			return
		}

		//Check in case we read 0 bytes
		if n > 0 {
			s.addRecBytes(n)
			c.MessageChan <- conn.NewMessage(c, n, b)
		}
	}
}
