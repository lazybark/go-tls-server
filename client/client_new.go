package client

import (
	"sync"

	"github.com/lazybark/go-helpers/semver"
	"github.com/lazybark/go-tls-server/conn"
)

// New creates new Client with specified config or default parameters.
func New(conf *Config) *Client {
	client := new(Client)
	client.errChan = make(chan error, 3) //nolint:gomnd // false alarm
	client.ClientDoneChan = make(chan bool)
	client.messageChan = make(chan *conn.Message, 10) //nolint:gomnd // false alarm
	client.mu = &sync.RWMutex{}
	client.ver = semver.Ver{ //nolint:exhaustruct // false alarm
		Major:       3, //nolint:gomnd // false alarm
		Minor:       2, //nolint:gomnd // false alarm
		Patch:       0,
		Stable:      false,
		ReleaseNote: "beta",
	}

	if conf == nil {
		conf = new(Config)
		// Dropping all stats is the default behaviour.
		conf.DropOldStats = true
	}

	// Default terminator is the newline.
	if conf.MessageTerminator == 0 {
		conf.MessageTerminator = '\n'
	}

	// Default buffer is 128 B.
	if conf.BufferSize == 0 {
		conf.BufferSize = 128
	}

	if conf.ErrorPrefix == "" {
		conf.ErrorPrefix = "TLS_CLIENT"
	}

	client.conf = conf

	return client
}
