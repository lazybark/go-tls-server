package server

import "errors"

//ErrMessageSizeLimit is returned after message length
//is equal or over server max message size directive
var ErrMessageSizeLimit = errors.New("message size limits reached")

//ErrReaderClosedByContext is returned after connection was closed
//by context cancelFunc during reading operation
var ErrReaderClosedByContext = errors.New("reader closed by context")

//ErrStreamClosed is returned after io.EOF is appeared in TLS stream
var ErrStreamClosed = errors.New("stream closed")
