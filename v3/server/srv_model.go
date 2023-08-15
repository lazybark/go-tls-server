package server

import (
	"crypto/tls"
	"net"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/lazybark/go-helpers/semver"
	"github.com/lazybark/go-tls-server/v3/conn"
)

var ver = semver.Ver{
	Major:       3,
	Minor:       1,
	Patch:       0,
	Stable:      false,
	ReleaseNote: "beta",
}

type Server struct {
	ver semver.Ver

	timeStart time.Time

	//host = hostname of the server
	host string

	//connPool is a map of pointers to connections.
	//
	//In this case pointers are used to increase code readability and number of ops
	//needed to change conn state
	connPool map[string]*conn.Connection

	//connPoolMutex controls connPool
	connPoolMutex sync.RWMutex

	//listener is the interface that listens for new connections
	listener net.Listener

	//tlsConfig points to tls listener config
	tlsConfig *tls.Config

	//sConfig points to server config
	sConfig *Config

	//ErrChan is the channel to send errors into external routine
	ErrChan chan error

	//ServerDoneChan is the channel to recieve server stopping command
	ServerDoneChan chan bool

	//ConnChan is the channel to notify external routine about new connection
	ConnChan chan *conn.Connection

	//stat keeps connections stat by date
	stat      map[string]Stat
	statMutex sync.RWMutex

	//statOverall keeps stats for all working time
	statOverall *Stat

	//resolver is the HTTP server used to serve stat data
	resolver *chi.Mux
}

// Version returns app version
func (s *Server) Version() semver.Ver { return s.ver }

// Version returns app version string
func (s *Server) VersionString() string { return s.ver.String() }
