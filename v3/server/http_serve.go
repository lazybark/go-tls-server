package server

import (
	"fmt"
	"log"
	"net/http"
)

func (s *Server) serveHTTP() {
	err := http.ListenAndServe(s.sConfig.HttpStatAddr, s.resolver)
	if err != nil {
		log.Fatal(fmt.Errorf("[Server][Listen] listening on %s is impossible: %w", s.sConfig.HttpStatAddr, err))
	}
}

func (s *Server) setHTTPRoutes() {
	s.resolver.Get("/api_version", func(w http.ResponseWriter, r *http.Request) {
		s.serveApiVersion(w)
	})

	s.resolver.Get("/stats", func(w http.ResponseWriter, r *http.Request) {
		s.serveStatistic(w)
	})

}

func (s *Server) serveApiVersion(w http.ResponseWriter) {
	_, err := w.Write([]byte(fmt.Sprintf(ResultString, s.Version())))
	if err != nil && !s.sConfig.SuppressErrors {
		s.ErrChan <- fmt.Errorf("[serveApiVersion] %w", err)
	}
}

func (s *Server) returnUnknownInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(fmt.Sprintf(ResultError, 500, "unknown_error_occurred")))
	if err != nil && !s.sConfig.SuppressErrors {
		s.ErrChan <- fmt.Errorf("[serveApiVersion] %w", err)
	}
}
