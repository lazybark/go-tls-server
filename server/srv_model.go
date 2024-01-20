package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/lazybark/go-helpers/semver"
	"github.com/lazybark/go-tls-server/conn"
)

type Server struct {
	ver semver.Ver

	isActive bool
	mu       *sync.Mutex

	timeStart time.Time

	// host = hostname of the server.
	host string

	// connPool is a map of pointers to connections.
	//
	// In this case pointers are used to increase code readability and number of ops
	// needed to change conn state.
	connPool map[string]*conn.Connection

	// connPoolMutex controls connPool.
	connPoolMutex sync.RWMutex

	// listener is the interface that listens for new connections.
	listener net.Listener

	// tlsConfig points to tls listener config.
	tlsConfig *tls.Config

	// sConfig points to server config.
	sConfig *Config

	// errChan is the channel to send errors into external routine.
	errChan chan error

	// serverDoneChan is the channel to receive server stopping command.
	serverDoneChan chan bool

	// connChan is the channel to notify external routine about new connection.
	connChan chan *conn.Connection

	// stat keeps connections stat by date.
	stat      map[string]Stat
	statMutex sync.RWMutex

	// statOverall keeps stats for all working time.
	statOverall *Stat

	// ctx is the server context.
	ctx    context.Context //nolint:containedctx // In TODOs
	cancel context.CancelFunc

	// errorPrefix is used as prefix to all errors to identify specific instance of server.
	//
	// Default: "TLS_SERVER".
	errorPrefix string
}

// Version returns app version.
func (s *Server) Version() semver.Ver { return s.ver }

// Version returns app version string.
func (s *Server) VersionString() string { return s.ver.String() }

// ErrChan returns server error channel available to read only.
func (s *Server) ErrChan() <-chan error { return s.errChan }

// ConnChan returns server new connections channel available to read only.
func (s *Server) ConnChan() <-chan *conn.Connection { return s.connChan }

// FormatError adds server's error prefix to err.
func (s *Server) FormatError(err error) error { return fmt.Errorf("%s: %w", s.errorPrefix, err) }

// SetActive sets server status to the value of active.
func (s *Server) SetActive(active bool) {
	s.mu.Lock()
	s.isActive = active
	s.mu.Unlock()
}

// IsActive returns current server status.
func (s *Server) IsActive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.isActive
}
