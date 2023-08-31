package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) serveStatistic(w http.ResponseWriter) {
	sentBytes, recievedBytes, errors, err := s.StatsOverall() //Error is always
	if err != nil && !s.sConfig.SuppressErrors {
		s.errChan <- fmt.Errorf("[Server][serveStatistic] %w", err)

		s.returnUnknownInternalError(w)
		return
	}

	conns, err := s.StatsConnections()
	if err != nil && !s.sConfig.SuppressErrors {
		s.errChan <- fmt.Errorf("[Server][serveStatistic] %w", err)

		s.returnUnknownInternalError(w)
		return
	}

	o := ServerStatsOutput{
		Sent:        sentBytes,
		Received:    recievedBytes,
		Errors:      errors,
		Connections: conns,
		Started:     s.timeStart,
	}

	oj, err := json.Marshal(o)
	if err != nil && !s.sConfig.SuppressErrors {
		s.errChan <- fmt.Errorf("[Server][serveStatistic] %w", err)

		s.returnUnknownInternalError(w)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf(ResultJSON, string(oj))))
	if err != nil && !s.sConfig.SuppressErrors {
		s.errChan <- fmt.Errorf("[serveApiVersion] %w", err)
	}
}
