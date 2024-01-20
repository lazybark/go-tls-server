package server

import (
	"fmt"
	"time"
)

// adminRoutine controls server behaviour: drops closed connections,
// closes inactive ones and stops the server in case s.ServerDoneChan.
func (s *Server) adminRoutine() { //nolint:cyclop,gocognit // in TODOs
	for {
		select {
		// Once per hour clean up old & close inactive connections.
		case <-time.After(time.Hour):
			for _, connection := range s.connPool {
				// If conn is closed and time now is already after the moment it should be deleted permanently.
				if connection.Closed() && !time.Now().
					Before(connection.ClosedAt().
						Add(time.Minute*time.Duration(s.sConfig.KeepOldConnections))) {
					s.remFromPool(connection)

					continue
				}

				// If it's not closed, but it's been a 'KeepInactiveConnections' time after.
				if !connection.Closed() && s.sConfig.KeepInactiveConnections > 0 &&
					!time.Now().Before(connection.LastAct().Add(time.Minute*time.Duration(s.sConfig.KeepInactiveConnections))) {
					err := s.CloseConnection(connection)
					if err != nil && !s.sConfig.SuppressErrors {
						s.errChan <- s.FormatError(fmt.Errorf("[Listen] error closing connection: %w", err))
					}
				}
			}
		// In case server needs to be stopped - close all connections.
		case d := <-s.serverDoneChan:
			if d {
				err := s.listener.Close()
				if err != nil && !s.sConfig.SuppressErrors {
					s.errChan <- s.FormatError(fmt.Errorf("[Listen] error closing listener: %w", err))
				}

				for _, c := range s.connPool {
					err := s.CloseConnection(c)
					if err != nil && !s.sConfig.SuppressErrors {
						s.errChan <- s.FormatError(fmt.Errorf("[adminRoutine] error closing connection %s -> %w", c.ID(), err))
					}
				}
			}
		}
	}
}
