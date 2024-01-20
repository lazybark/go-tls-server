package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"

	"github.com/lazybark/go-tls-server/conn"
)

// Listen runs listener interface implementations and accepts connections.
func (s *Server) Listen(port string) error { //nolint:cyclop // in TODOs
	listener, err := tls.Listen("tcp", ":"+port, s.tlsConfig)
	if err != nil {
		return s.FormatError(fmt.Errorf("[Listen] error listening: %w", err))
	}

	s.SetActive(true)

	s.listener = listener

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
				// Accept the connection.
				tlsConn, err := listener.Accept()

				// The problem is that a listener can be closed during the listening. Then we get net.ErrClosed.
				// In this case we always ignore it, because doesn't matter why it's closed: this function is not for err processing.
				// Error was handled somewhere else already. Or server was simply terminated.
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						continue
					}

					if !s.sConfig.SuppressErrors {
						s.errChan <- s.FormatError(fmt.Errorf("[Listen] error accepting connection: %w", err))
					}
				}

				// Just a precaution to avoid nil pointer dereference.
				if tlsConn == nil {
					continue
				}

				connection, err := conn.NewConnection(tlsConn.RemoteAddr(), tlsConn, s.sConfig.MessageTerminator)
				if err != nil && !s.sConfig.SuppressErrors {
					s.errChan <- s.FormatError(fmt.Errorf("[Listen] error making connection for %v: %w", tlsConn.RemoteAddr(), err))
				}

				// Add to pool.
				s.addToPool(connection)
				// Notify outer routine.
				s.connChan <- connection
				// Wait for new messages.
				go s.receive(connection)
			}
		}
	}()

	return nil
}

// receive endlessly reads incoming stream and delivers messages to receivers outside server routine.
// It uses ReadWithContext, so execution can be manually stopped by calling c.cancel on specific connection.
// In that case (or if any error occurs) method will trigger s.CloseConnection to break connection too.
func (s *Server) receive(connection *conn.Connection) {
	for {
		if connection.Closed() {
			return
		}

		bytes, bytesCount, err := connection.ReadWithContext(
			s.sConfig.BufferSize,
			s.sConfig.MaxMessageSize,
			s.sConfig.MessageTerminator,
		)
		if err != nil {
			if !s.sConfig.SuppressErrors {
				s.errChan <- s.FormatError(fmt.Errorf("[receive] error reading from %s: %w", connection.ID(), err))
			}

			err := s.CloseConnection(connection)
			if err != nil && !s.sConfig.SuppressErrors {
				s.errChan <- s.FormatError(fmt.Errorf("[receive] error closing connection: %w", err))
			}

			return
		}

		// Check in case we read 0 bytes.
		if bytesCount > 0 {
			s.addRecBytes(bytesCount)
			connection.MessageChanWrite() <- conn.NewMessage(connection, bytesCount, bytes)
		}
	}
}
