package client

type Config struct {
	// SuppressErrors prevents client from sending errors into ErrChan.
	// Does not include fatal errors during startup.
	SuppressErrors bool

	// MaxMessageSize sets max length of one message in bytes.
	// If >0 and limit is reached, connection will be closed with an error.
	//
	// Note that if MaxMessageSize is > than reading buffer and MaxMessageSize reached,
	// it will not close connection until buffer is full or message terminator occurs.
	MaxMessageSize int

	// MessageTerminator sets byte value that marks end of the message in stream.
	// Works for both incoming and outgoing messages.
	MessageTerminator byte

	// BufferSize regulates buffer length to read incoming message. Default value is 128.
	BufferSize int

	// DropOldStats = true will make client to set all sent/received bytes & errors to zero before opening new connection.
	DropOldStats bool

	// errorPrefix is used as prefix to all errors to identify specific instance of client.
	//
	// Default: "TLS_CLIENT".
	ErrorPrefix string
}
