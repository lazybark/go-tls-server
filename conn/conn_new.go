package conn

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/lazybark/go-helpers/npt"
)

func NewConnection(address net.Addr, conn net.Conn, terminator byte) (*Connection, error) {
	connection := new(Connection)
	connection.connectedAt = npt.Now()
	connection.lastAct = connection.connectedAt
	connection.messageChan = make(chan *Message)
	connection.mu = &sync.RWMutex{}

	connID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("[NewConnection] %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	connection.id = connID.String()
	connection.addr = address
	connection.tlsConn = conn
	connection.cancel = cancel
	connection.ctx = ctx
	connection.SetMessageTerminator(terminator)

	return connection, nil
}
