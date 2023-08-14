package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) serveStatistic(w http.ResponseWriter) {
	sentBytes, recievedBytes, errors, err := s.StatsOverall() //Error is always
	if err != nil && !s.sConfig.SuppressErrors {
		s.ErrChan <- fmt.Errorf("[Server][serveStatistic] %w", err)

		returnUnknownInternalError(w)
		return
	}

	conns, err := s.StatsConnections()
	if err != nil && !s.sConfig.SuppressErrors {
		s.ErrChan <- fmt.Errorf("[Server][serveStatistic] %w", err)

		returnUnknownInternalError(w)
		return
	}

	o := ServerStatsOutput{
		Sent:        sentBytes,
		Recieved:    recievedBytes,
		Errors:      errors,
		Connections: conns,
	}

	oj, err := json.Marshal(o)
	if err != nil && !s.sConfig.SuppressErrors {
		s.ErrChan <- fmt.Errorf("[Server][serveStatistic] %w", err)

		returnUnknownInternalError(w)
		return
	}

	w.Write([]byte(fmt.Sprintf(ResultJSON, string(oj))))
}
