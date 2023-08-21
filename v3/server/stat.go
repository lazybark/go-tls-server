package server

import (
	"fmt"
	"time"
)

// Stat holds server statistic
type Stat struct {
	received int
	sent     int
	errors   int
}

var ErrNoStatForTheDay = fmt.Errorf("no stat")

const statKeyPattern string = "%d-%d-%d"

func getStatKey() string {
	now := time.Now()
	return fmt.Sprintf(statKeyPattern, now.Year(), now.Month(), now.Day())
}

// Stats returns server stats for specified day or error in case date is not in stat
func (s *Server) StatsOverall() (sentBytes, receivedBytes, errors int, err error) {
	//Right now err is not used - added for compatibility for future

	s.statMutex.Lock()
	defer s.statMutex.Unlock()

	sentBytes = s.statOverall.sent
	receivedBytes = s.statOverall.received
	errors = s.statOverall.errors

	return
}

// Stats returns server stats for specified day or error in case date is not in stat
func (s *Server) Stats(y int, m int, d int) (sentBytes, receivedBytes, errors int, err error) {
	date := fmt.Sprintf(statKeyPattern, y, m, d)

	s.statMutex.Lock()
	defer s.statMutex.Unlock()
	if v, ok := s.stat[date]; ok {
		return v.sent, v.received, v.errors, nil
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
	if v, ok := s.stat[d]; ok {
		v.received += n
		s.stat[d] = v
	} else {
		s.stat[d] = Stat{received: n}
	}

	s.statOverall.received += n

}

// addSentBytes adds bytes to stat of current day
func (s *Server) addSentBytes(n int) {
	if n < 0 {
		return
	}
	d := getStatKey()

	s.statMutex.Lock()
	defer s.statMutex.Unlock()
	if v, ok := s.stat[d]; ok {
		v.sent += n
		s.stat[d] = v
	} else {
		s.stat[d] = Stat{sent: n}
	}

	s.statOverall.sent += n
}

// addErrors adds bytes to stat of current day
func (s *Server) addErrors(n int) {
	if n < 0 {
		return
	}
	d := getStatKey()

	s.statMutex.Lock()
	defer s.statMutex.Unlock()
	if v, ok := s.stat[d]; ok {
		v.errors += n
		s.stat[d] = v
	} else {
		s.stat[d] = Stat{errors: n}
	}

	s.statOverall.errors += n
}

// StartedAt returns starting time
func (s *Server) StartedAt() time.Time { return s.timeStart }

// Online returns time online
func (s *Server) Online() time.Duration { return time.Since(s.timeStart) }

// ActiveConnetions returns number of active connections
func (s *Server) ActiveConnetions() int {
	a := 0
	for _, c := range s.connPool {
		if !c.Closed() {
			a++
		}
	}
	return a
}
