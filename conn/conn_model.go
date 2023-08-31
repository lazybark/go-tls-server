package conn

import (
	"context"
	"net"
	"sync"

	"github.com/lazybark/go-helpers/npt"
)

// Connection represents incoming client connection
type Connection struct {
	//id is an unique string to define connection
	id string

	//connectedAt time of connection init
	connectedAt npt.NPT

	//addr is the remote address of client
	addr net.Addr

	//conn is the connection interface that reads and writes bytes
	tlsConn net.Conn

	//isClosed = true means that connection was closed and soon will be dropped from pool
	isClosed bool

	//closedAt is the time connection was marked as 'closed'
	closedAt npt.NPT

	//lastAct updates every time there was any action in connection
	lastAct npt.NPT

	//ctx is the connection context
	ctx    context.Context
	cancel context.CancelFunc

	//bytesLeft holds extra bytes that were read from stream after terminator occurred, but end of buffer was not reached
	bytesLeft []byte

	//bs holds total bytes sent by server in connection
	bs int

	//br holds total bytes received by server in connection
	br int

	//errors holds total number of errors occurred in connection
	errors int

	//MessageTerminator sets byte value that marks message end in the stream.
	//Works for both incoming and outgoing messages
	messageTerminator byte

	//messageChan channel to notify external routine about new messages
	messageChan chan *Message

	mu *sync.RWMutex
}

// Address returns remote address of client
func (c *Connection) Address() net.Addr { return c.addr }

// Id returns connection ID in pool
func (c *Connection) Id() string { return c.id }

// SetMessageTerminator sets byte that will be used as message terminator
func (c *Connection) SetMessageTerminator(t byte) { c.messageTerminator = t }

// MessageChanRead returns connection's message channel to read only
func (c *Connection) MessageChanRead() <-chan *Message { return c.messageChan }

// MessageChanWrite returns connection's message channel to write only
func (c *Connection) MessageChanWrite() chan<- *Message { return c.messageChan }
