package conn

import (
	"context"
	"net"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/lazybark/go-helpers/npt"
)

func NewConnection(ip net.Addr, conn net.Conn, t byte) (*Connection, error) {
	c := new(Connection)
	c.connectedAt = npt.Now()
	c.lastAct = c.connectedAt
	c.messageChan = make(chan *Message)
	c.mu = &sync.RWMutex{}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())

	c.id = id.String()
	c.addr = ip
	c.tlsConn = conn
	c.cancel = cancel
	c.ctx = ctx
	c.SetMessageTerminator(t)

	return c, nil
}
