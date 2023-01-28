package server

//Message represents incoming message with its bytes and pointer to associated connection
type Message struct {
	conn   *Connection
	length int
	bytes  []byte
}

//Bytes returns message bytes
func (m *Message) Bytes() []byte { return m.bytes }

//Length returns message bytes length
func (m *Message) Length() int { return m.length }

//Length returns message bytes length
func (m *Message) Conn() *Connection { return m.conn }
