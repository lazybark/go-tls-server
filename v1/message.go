package v1

//Message represents incoming message with its bytes and pointer to associated connection
type Message struct {
	conn  *Connection
	bytes []byte
}

//Bytes returns message bytes
func (m *Message) Bytes() []byte {
	return m.bytes
}
