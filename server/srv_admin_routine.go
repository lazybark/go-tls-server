package server

import (
	"fmt"
	"time"
)

// adminRoutine controls server behaviour: drops closed connections, closes inactive ones and stops the server in case s.ServerDoneChan
func (s *Server) adminRoutine() {
	for {
		select {
		//Once per hour clean up old & close inactive connections
		case <-time.After(time.Hour):
			for _, c := range s.connPool {
				//If conn is closed and time now is already after the moment it should be deleted permanently
				if c.Closed() && !time.Now().Before(c.ClosedAt().Add(time.Minute*time.Duration(s.sConfig.KeepOldConnections))) {
					s.remFromPool(c)
					continue
				}
				//If it's not closed, but it's been a 'KeepInactiveConnections' time after
				if !c.Closed() && s.sConfig.KeepInactiveConnections > 0 && !time.Now().Before(c.LastAct().Add(time.Minute*time.Duration(s.sConfig.KeepInactiveConnections))) {
					err := s.CloseConnection(c)
					if err != nil && !s.sConfig.SuppressErrors {
						s.errChan <- fmt.Errorf("[Server][Listen] error closing connection: %w", err)
					}
				}
			}
		//In case server needs to be stopped - close all connections
		case d := <-s.serverDoneChan:
			if d {
				err := s.listener.Close()
				if err != nil && !s.sConfig.SuppressErrors {
					s.errChan <- fmt.Errorf("[Server][Listen] error closing listener: %w", err)
				}
				for _, c := range s.connPool {
					err := s.CloseConnection(c)
					if err != nil && !s.sConfig.SuppressErrors {
						s.errChan <- fmt.Errorf("[Server][adminRoutine] error closing connection %s -> %w", c.Id(), err)
					}
				}
			}
		}
	}
}
