package conn

import "errors"

//ErrMessageSizeLimit is returned after message length
//is over server max message size directive
var ErrMessageSizeLimit = errors.New("message size limits reached")

//ErrReaderClosedByContext is returned after connection was closed
//by context cancelFunc during reading operation
var ErrReaderClosedByContext = errors.New("reader closed by context")

//ErrReaderAlreadyClosed is returned when client code attempts to read from
//previously closed connection
var ErrReaderAlreadyClosed = errors.New("reader already closed")

//ErrStreamClosed is returned after io.EOF is appeared in TLS stream
var ErrStreamClosed = errors.New("stream closed")
