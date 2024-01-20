package server

type Config struct {
	// SuppressErrors prevents server from sending errors into ErrChan.
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

	// KeepOldConnections prevents server from dropping closed connection for N minutes after it has been closed.
	// Useful for keeping stats, but it's deadly to keep them forever.
	//
	// Default value (in case 0) will be set to 1440 min (24h).
	KeepOldConnections int

	// KeepInactiveConnections makes server close connection that had no activity for N mins.
	// 0 means keep such connection forever.
	KeepInactiveConnections int

	// ErrorPrefix is used as prefix to all errors to identify specific instance of server.
	//
	// Default: "TLS_SERVER"
	ErrorPrefix string
}
