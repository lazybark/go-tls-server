package server

import (
	"github.com/lazybark/go-tls-server/conn"
)

// CloseConnection is the only correct way to close connection.
// It changes conn state in pool and then calls to c.close.
func (s *Server) CloseConnection(c *conn.Connection) error {
	err := c.Close()
	if err != nil {
		return s.FormatError(err)
	}

	return nil
}

// addToPool adds connection to fool for controlling.
func (s *Server) addToPool(c *conn.Connection) {
	s.connPoolMutex.Lock()
	s.connPool[c.ID()] = c
	s.connPoolMutex.Unlock()
}

// remFromPool removes connection pointer from pool, so it becomes unavailable to reach.
func (s *Server) remFromPool(c *conn.Connection) {
	s.connPoolMutex.Lock()
	delete(s.connPool, c.ID())
	s.connPoolMutex.Unlock()
}

// SendByte calls to c.SendByte and adds sent bytes to Stat.
func (s *Server) SendByte(c *conn.Connection, b []byte) error {
	n, err := c.SendByte(b)

	s.addSentBytes(n)

	if err != nil {
		s.addErrors(1)
	}

	return s.FormatError(err)
}

// SendString calls to c.SendString and adds sent bytes to Stat.
func (s *Server) SendString(c *conn.Connection, str string) error {
	n, err := c.SendString(str)
	if err != nil {
		s.addErrors(1)
	}

	s.addSentBytes(n)

	return s.FormatError(err)
}
