package server

import (
	"errors"
	"fmt"
	"time"
)

// Stat holds server statistic.
type Stat struct {
	received int
	sent     int
	errors   int
}

var ErrNoStatForTheDay = errors.New("no stat")

const statKeyPattern string = "%d-%d-%d"

func getStatKey() string {
	now := time.Now()

	return fmt.Sprintf(statKeyPattern, now.Year(), now.Month(), now.Day())
}

// Stats returns server stats for specified day or error in case date is not in stat.
//
// Right now err is not used - added for compatibility for future.
func (s *Server) StatsOverall() (int, int, int, error) {
	s.statMutex.Lock()
	defer s.statMutex.Unlock()

	return s.statOverall.sent, s.statOverall.received, s.statOverall.errors, nil
}

// Stats returns server stats for specified day or error in case date is not in stat.
func (s *Server) Stats(y int, m int, d int) (int, int, int, error) {
	date := fmt.Sprintf(statKeyPattern, y, m, d)

	s.statMutex.Lock()
	defer s.statMutex.Unlock()

	if v, ok := s.stat[date]; ok {
		return v.sent, v.received, v.errors, nil
	}

	return 0, 0, 0, ErrNoStatForTheDay
}

// StatsConnections returns statistic about current connections.
//
// Right now err is not used - added for compatibility for future.
func (s *Server) StatsConnections() (int, error) {
	s.connPoolMutex.Lock()
	defer s.connPoolMutex.Unlock()

	return len(s.connPool), nil
}

// addRecBytes adds bytes to stat of current day.
func (s *Server) addRecBytes(count int) {
	if count < 0 {
		return
	}

	date := getStatKey()

	s.statMutex.Lock()
	defer s.statMutex.Unlock()

	if v, ok := s.stat[date]; ok {
		v.received += count
		s.stat[date] = v
	} else {
		s.stat[date] = Stat{received: count, sent: 0, errors: 0}
	}

	s.statOverall.received += count
}

// addSentBytes adds bytes to stat of current day.
func (s *Server) addSentBytes(count int) {
	if count < 0 {
		return
	}

	date := getStatKey()

	s.statMutex.Lock()
	defer s.statMutex.Unlock()

	if v, ok := s.stat[date]; ok {
		v.sent += count
		s.stat[date] = v
	} else {
		s.stat[date] = Stat{sent: count, received: 0, errors: 0}
	}

	s.statOverall.sent += count
}

// addErrors adds bytes to stat of current day.
func (s *Server) addErrors(count int) {
	if count < 0 {
		return
	}

	date := getStatKey()

	s.statMutex.Lock()
	defer s.statMutex.Unlock()

	if v, ok := s.stat[date]; ok {
		v.errors += count
		s.stat[date] = v
	} else {
		s.stat[date] = Stat{errors: count, sent: 0, received: 0}
	}

	s.statOverall.errors += count
}

// StartedAt returns starting time.
func (s *Server) StartedAt() time.Time { return s.timeStart }

// Online returns time online.
func (s *Server) Online() time.Duration { return time.Since(s.timeStart) }

// ActiveConnetions returns number of active connections.
func (s *Server) ActiveConnetions() int {
	connection := 0

	for _, c := range s.connPool {
		if !c.Closed() {
			connection++
		}
	}

	return connection
}
