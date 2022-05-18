package v1

import (
	"fmt"
	"time"
)

//Stat holds server statistic
type Stat struct {
	recieved    int
	sent        int
	errors      int
	connections int
}

//Stats returns server stats for specified day or error in case date is not in stat
func (s *Server) Stats(y int, m int, d int) (int, int, int, int, error) {
	date := fmt.Sprintf("%d-%d-%d", y, m, d)

	if v, ok := s.Stat[date]; ok {
		return v.sent, v.recieved, v.errors, v.connections, nil
	}
	return 0, 0, 0, 0, fmt.Errorf("no stat")
}

//addRecBytes adds bytes to stat of current day
func (s *Server) addRecBytes(n int) {
	d := fmt.Sprintf("%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day())
	if v, ok := s.Stat[d]; ok {
		v.recieved += n
		s.Stat[d] = v
	} else {
		s.Stat[d] = Stat{recieved: n}
	}
}

//addSentBytes adds bytes to stat of current day
func (s *Server) addSentBytes(n int) {
	d := fmt.Sprintf("%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day())
	if v, ok := s.Stat[d]; ok {
		v.sent += n
		s.Stat[d] = v
	} else {
		s.Stat[d] = Stat{sent: n}
	}
}

//addErrors adds bytes to stat of current day
func (s *Server) addErrors(n int) {
	d := fmt.Sprintf("%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day())
	if v, ok := s.Stat[d]; ok {
		v.errors += n
		s.Stat[d] = v
	} else {
		s.Stat[d] = Stat{errors: n}
	}
}
