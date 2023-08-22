package conn

// Message represents incoming message with its bytes and pointer to associated connection
type Message struct {
	conn   *Connection
	length int
	bytes  []byte
}

func NewMessage(conn *Connection, length int, bytes []byte) *Message {
	return &Message{conn: conn, length: length, bytes: bytes}
}

// Bytes returns message bytes
func (m *Message) Bytes() []byte { return m.bytes }

// Length returns message bytes length
func (m *Message) Length() int { return m.length }

// Conn returns pointer to connection in which message was received
func (m *Message) Conn() *Connection { return m.conn }
