package conn

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testMessageTerminator = "\n"
)

func TestConnectionCorrectByteSending(t *testing.T) {
	send := "Hello there, General Kenobi!"
	tlsConn := &MockTLSConnection{}

	cn, err := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	cn.SendByte([]byte(send))
	assert.Equal(t, send+testMessageTerminator, string(tlsConn.WriteBytesTo))

	sent, rec, errs := cn.Stats()
	assert.Equal(t, len(send+testMessageTerminator), sent)
	assert.Equal(t, 0, errs)
	assert.Equal(t, 0, rec)
}

func TestConnectionCorrectStringSending(t *testing.T) {
	send := "Hello there, General Kenobi!"
	tlsConn := &MockTLSConnection{}

	cn, err := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	cn.SendString(send)
	assert.Equal(t, send+testMessageTerminator, string(tlsConn.WriteBytesTo))

	sent, rec, errs := cn.Stats()
	assert.Equal(t, len(send+testMessageTerminator), sent)
	assert.Equal(t, 0, errs)
	assert.Equal(t, 0, rec)
}

func TestConnectionCorrectReading(t *testing.T) {
	tlsConn := &MockTLSConnection{
		BytesToRead: []byte("Hello there, General Kenobi!\n"),
	}

	cn, err := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	buffer := 5
	maxSize := 550
	read, count, err := cn.ReadWithContext(buffer, maxSize, testMessageTerminator[0])
	require.NoError(t, err)

	assert.Equal(t, string(tlsConn.BytesToRead[:len(tlsConn.BytesToRead)-1]), string(read))
	assert.Equal(t, len(tlsConn.BytesToRead), count)

	sent, rec, errs := cn.Stats()
	assert.Equal(t, len(tlsConn.BytesToRead), rec)
	assert.Equal(t, 0, errs)
	assert.Equal(t, 0, sent)

}

func TestConnectionReadingTooLargeMessage(t *testing.T) {
	tlsConn := &MockTLSConnection{
		BytesToRead: []byte("Hello there, General Kenobi!\n"),
	}

	cn, err := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	//Small buffer will result in double buffer reads until max can be checked
	buffer := 5
	maxSize := 5
	read, count, err := cn.ReadWithContext(buffer, maxSize, testMessageTerminator[0])
	assert.Equal(t, true, errors.Is(err, ErrMessageSizeLimit))

	assert.Equal(t, "", string(read)) //Returning nothing
	assert.Equal(t, buffer*2, count)  //Reading performed two times before got more than allowed

	sent, rec, errs := cn.Stats()
	assert.Equal(t, buffer*2, rec) //Reading performed two times before got more than allowed
	assert.Equal(t, 1, errs)       //Got exactly 1 error
	assert.Equal(t, 0, sent)

}

func TestConnectionReadingTooLargeMessageTooLargeBuffer(t *testing.T) {
	tlsConn := &MockTLSConnection{
		BytesToRead: []byte("Hello there, General Kenobi!\n"),
	}

	cn, err := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	//buffer > message length will result in error after whole message read
	buffer := 50
	maxSize := 5
	read, count, err := cn.ReadWithContext(buffer, maxSize, testMessageTerminator[0])
	assert.Equal(t, true, errors.Is(err, ErrMessageSizeLimit))

	assert.Equal(t, "", string(read))                //Returning nothing
	assert.Equal(t, len(tlsConn.BytesToRead), count) //Larger buffer reads whole message before there is an error

	sent, rec, errs := cn.Stats()
	assert.Equal(t, len(tlsConn.BytesToRead), rec) //Larger buffer reads whole message before there is an error
	assert.Equal(t, 1, errs)                       //Got exactly 1 error
	assert.Equal(t, 0, sent)

}

func TestConnectionReadingClosedByContext(t *testing.T) {
	tlsConn := &MockTLSConnection{
		BytesToRead: []byte("Hello there, General Kenobi!"), //No terminator, so it will go until stopped
	}

	cn, err := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	buffer := 5
	maxSize := 100
	var read []byte
	var count int

	go func(read *[]byte, count *int, err *error) {
		*read, *count, *err = cn.ReadWithContext(buffer, maxSize, testMessageTerminator[0])
	}(&read, &count, &err)

	time.Sleep(time.Second * 3)
	cn.cancel()

	assert.Equal(t, nil, err)

	assert.Equal(t, len(tlsConn.BytesToRead), count) //It still returns count of bytes read
	assert.Equal(t, "", string(read))                //But not the data

	sent, rec, errs := cn.Stats()
	assert.Equal(t, len(tlsConn.BytesToRead), rec) //It still returns count of bytes read
	assert.Equal(t, 0, errs)                       //Errors 0, as there was no error itself, just stopped by externall call
	assert.Equal(t, 0, sent)
}

func TestConnectionReadingClose(t *testing.T) {
	tlsConn := &MockTLSConnection{
		BytesToRead: []byte("Hello there, General Kenobi!"), //No terminator, so it will go until stopped
	}

	cn, err := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	buffer := 5
	maxSize := 100
	var read []byte
	var count int

	go func(read *[]byte, count *int, err *error) {
		*read, *count, *err = cn.ReadWithContext(buffer, maxSize, testMessageTerminator[0])
	}(&read, &count, &err)

	time.Sleep(time.Second * 3)
	cn.Close()

	assert.Equal(t, nil, err)

	assert.Equal(t, len(tlsConn.BytesToRead), count) //It still returns count of bytes read
	assert.Equal(t, "", string(read))                //But not the data

	sent, rec, errs := cn.Stats()
	assert.Equal(t, len(tlsConn.BytesToRead), rec) //It still returns count of bytes read
	assert.Equal(t, 0, errs)                       //Errors 0, as there was no error itself, just stopped by externall call
	assert.Equal(t, 0, sent)

	assert.Equal(t, true, cn.isClosed)
	assert.Equal(t, true, cn.Closed())
	assert.Equal(t, true, tlsConn.AskedToBeClosed)
}

func TestConnectionReadingEOF(t *testing.T) {
	tlsConn := &MockTLSConnection{
		BytesToRead: []byte("Hello there, General Kenobi!"), //No terminator, so it will go until stopped
		ReturnEOF:   true,
	}

	cn, err := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	buffer := 50
	maxSize := 5

	read, count, err := cn.ReadWithContext(buffer, maxSize, testMessageTerminator[0])
	assert.Equal(t, true, errors.Is(err, ErrStreamClosed))

	assert.Equal(t, "", string(read))
	assert.Equal(t, 0, count)

	sent, rec, errs := cn.Stats()
	assert.Equal(t, 0, rec)  //Larger buffer reads whole message before there is an error
	assert.Equal(t, 0, errs) //Got exactly 1 error
	assert.Equal(t, 0, sent)
}

func TestConnectionReadingStats(t *testing.T) {
	tlsConn := &MockTLSConnection{
		BytesToRead: []byte("Hello there, General Kenobi!"), //No terminator, so it will go until stopped
	}

	cn, err := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	cn.br = 3333
	cn.errors = 15
	cn.bs = 8374

	sent, rec, errs := cn.Stats()
	assert.Equal(t, cn.br, rec)
	assert.Equal(t, cn.errors, errs)
	assert.Equal(t, cn.bs, sent)
	assert.Equal(t, cn.br, 3333)
	assert.Equal(t, cn.errors, 15)
	assert.Equal(t, cn.bs, 8374)

	cn.DropOldStats()

	sent, rec, errs = cn.Stats()
	assert.Equal(t, 0, rec)
	assert.Equal(t, 0, errs)
	assert.Equal(t, 0, sent)

}
