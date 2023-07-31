package server

import (
	"fmt"
	"time"
)

// Stat holds server statistic
type Stat struct {
	recieved int
	sent     int
	errors   int
}

var ErrNoStatForTheDay = fmt.Errorf("no stat")

const statKeyPattern = "%d-%d-%d"

// Stats returns server stats for specified day or error in case date is not in stat
func (s *Server) Stats(y int, m int, d int) (sentBytes, recievedBytes, errors int, err error) {
	date := fmt.Sprintf(statKeyPattern, y, m, d)

	s.statMutex.Lock()
	defer s.statMutex.Unlock()
	if v, ok := s.Stat[date]; ok {
		return v.sent, v.recieved, v.errors, nil
	}

	err = ErrNoStatForTheDay
	return
}

// StatsConnections returns statistic about current connections
func (s *Server) StatsConnections() (connections int, err error) {
	//Right now err is not used - added for compatibility for future
	s.connPoolMutex.Lock()
	connections = len(s.connPool)
	s.connPoolMutex.Unlock()

	return
}

// addRecBytes adds bytes to stat of current day
func (s *Server) addRecBytes(n int) {
	if n < 0 {
		return
	}
	d := getStatKey()

	s.statMutex.Lock()
	defer s.statMutex.Unlock()
	if v, ok := s.Stat[d]; ok {
		v.recieved += n
		s.Stat[d] = v
	} else {
		s.Stat[d] = Stat{recieved: n}
	}

}

// addSentBytes adds bytes to stat of current day
func (s *Server) addSentBytes(n int) {
	if n < 0 {
		return
	}
	d := getStatKey()

	s.statMutex.Lock()
	defer s.statMutex.Unlock()
	if v, ok := s.Stat[d]; ok {
		v.sent += n
		s.Stat[d] = v
	} else {
		s.Stat[d] = Stat{sent: n}
	}
}

// addErrors adds bytes to stat of current day
func (s *Server) addErrors(n int) {
	if n < 0 {
		return
	}
	d := getStatKey()

	s.statMutex.Lock()
	defer s.statMutex.Unlock()
	if v, ok := s.Stat[d]; ok {
		v.errors += n
		s.Stat[d] = v
	} else {
		s.Stat[d] = Stat{errors: n}
	}
}

func getStatKey() string {
	now := time.Now()
	return fmt.Sprintf(statKeyPattern, now.Year(), now.Month(), now.Day())
}
