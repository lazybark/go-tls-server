package conn_test

import (
	"errors"
	"testing"
	"time"

	"github.com/lazybark/go-helpers/mock"
	"github.com/lazybark/go-tls-server/conn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testMessageTerminator = "\n"
)

func TestConnectionCorrectByteSending(t *testing.T) {
	send := "Hello there, General Kenobi!"
	tlsConn := &mock.MockTLSConnection{}

	cn, err := conn.NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	_, _ = cn.SendByte([]byte(send))
	assert.Equal(t, send+testMessageTerminator, string(tlsConn.MWR.Bytes))

	sent, rec, errs := cn.Stats()
	assert.Equal(t, len(send+testMessageTerminator), sent) //Sending bytes, so counting also bytes, not chars
	assert.Equal(t, 0, errs)
	assert.Equal(t, 0, rec)
}

func TestConnectionCorrectStringSending(t *testing.T) {
	send := "Hello there, General Kenobi!"
	tlsConn := &mock.MockTLSConnection{}

	cn, err := conn.NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	_, _ = cn.SendString(send)
	assert.Equal(t, send+testMessageTerminator, string(tlsConn.MWR.Bytes))

	sent, rec, errs := cn.Stats()
	assert.Equal(t, len(send+testMessageTerminator), sent) //Sending bytes, so counting also bytes, not chars
	assert.Equal(t, 0, errs)
	assert.Equal(t, 0, rec)
}

func TestConnectionCorrectReading(t *testing.T) {
	tlsConn := &mock.MockTLSConnection{
		MWR: mock.MockWriteReader{
			Bytes:            []byte("Hello there, General Kenobi!\n"),
			DontReturEOFEver: true, //To mock endless stream with opened client
		},
	}

	cn, err := conn.NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	buffer := 5
	maxSize := 550
	read, count, err := cn.ReadWithContext(buffer, maxSize, testMessageTerminator[0])
	require.NoError(t, err)

	assert.Equal(t, string(tlsConn.MWR.Bytes[:len(tlsConn.MWR.Bytes)-1]), string(read))
	assert.Equal(t, len(tlsConn.MWR.Bytes), count)

	sent, rec, errs := cn.Stats()
	assert.Equal(t, len(tlsConn.MWR.Bytes), rec)
	assert.Equal(t, 0, errs)
	assert.Equal(t, 0, sent)

}

func TestConnectionReadingTooLargeMessage(t *testing.T) {
	tlsConn := &mock.MockTLSConnection{
		MWR: mock.MockWriteReader{
			Bytes:            []byte("Hello there, General Kenobi!\n"),
			DontReturEOFEver: true, //To mock endless stream with opened client
		},
	}

	cn, err := conn.NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	//Small buffer will result in double buffer reads until max can be checked
	buffer := 5
	maxSize := 5
	read, count, err := cn.ReadWithContext(buffer, maxSize, testMessageTerminator[0])
	assert.Equal(t, true, errors.Is(err, conn.ErrMessageSizeLimit))

	assert.Equal(t, "", string(read)) //Returning nothing
	assert.Equal(t, buffer*2, count)  //Reading performed two times before got more than allowed

	sent, rec, errs := cn.Stats()
	assert.Equal(t, buffer*2, rec) //Reading performed two times before got more than allowed
	assert.Equal(t, 1, errs)       //Got exactly 1 error
	assert.Equal(t, 0, sent)

}

func TestConnectionReadingTooLargeMessageTooLargeBuffer(t *testing.T) {
	tlsConn := &mock.MockTLSConnection{
		MWR: mock.MockWriteReader{
			Bytes:            []byte("Hello there, General Kenobi!\n"),
			DontReturEOFEver: true, //To mock endless stream with opened client
		},
	}

	cn, err := conn.NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	//buffer > message length will result in error after whole message read
	buffer := 50
	maxSize := 5
	read, count, err := cn.ReadWithContext(buffer, maxSize, testMessageTerminator[0])
	assert.Equal(t, true, errors.Is(err, conn.ErrMessageSizeLimit))

	assert.Equal(t, "", string(read))              //Returning nothing
	assert.Equal(t, len(tlsConn.MWR.Bytes), count) //Larger buffer reads whole message before there is an error

	sent, rec, errs := cn.Stats()
	assert.Equal(t, len(tlsConn.MWR.Bytes), rec) //Larger buffer reads whole message before there is an error
	assert.Equal(t, 1, errs)                     //Got exactly 1 error
	assert.Equal(t, 0, sent)

}

func TestConnectionReadingClosedByContext(t *testing.T) {
	tlsConn := &mock.MockTLSConnection{
		MWR: mock.MockWriteReader{
			Bytes:            []byte("Hello there, General Kenobi!"), //No terminator, so it will go until stopped
			DontReturEOFEver: true,                                   //To mock endless stream with opened client
		},
	}

	cn, err := conn.NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	buffer := 5
	maxSize := 100
	var read []byte
	var count int

	go func() {
		time.Sleep(time.Second * 3)
		cn.CancelCtx()
	}()

	read, count, err = cn.ReadWithContext(buffer, maxSize, testMessageTerminator[0])
	assert.Equal(t, nil, err)

	assert.Equal(t, len(tlsConn.MWR.Bytes), count) //It still returns count of bytes read
	assert.Equal(t, "", string(read))              //But not the data

	sent, rec, errs := cn.Stats()
	assert.Equal(t, len(tlsConn.MWR.Bytes), rec) //It still returns count of bytes read
	assert.Equal(t, 0, errs)                     //Errors 0, as there was no error itself, just stopped by externall call
	assert.Equal(t, 0, sent)
}

func TestConnectionReadingClose(t *testing.T) {
	tlsConn := &mock.MockTLSConnection{
		MWR: mock.MockWriteReader{
			Bytes:            []byte("Hello there, General Kenobi!"), //No terminator, so it will go until stopped
			DontReturEOFEver: true,                                   //To mock endless stream with opened client
		},
	}

	cn, err := conn.NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	buffer := 5
	maxSize := 100
	var read []byte
	var count int

	go func() {
		time.Sleep(time.Second * 3)
		_ = cn.Close()
	}()

	read, count, err = cn.ReadWithContext(buffer, maxSize, testMessageTerminator[0])
	assert.Equal(t, nil, err)

	assert.Equal(t, len(tlsConn.MWR.Bytes), count) //It still returns count of bytes read
	assert.Equal(t, "", string(read))              //But not the data

	sent, rec, errs := cn.Stats()
	assert.Equal(t, len(tlsConn.MWR.Bytes), rec) //It still returns count of bytes read
	assert.Equal(t, 0, errs)                     //Errors 0, as there was no error itself, just stopped by externall call
	assert.Equal(t, 0, sent)

	assert.Equal(t, true, cn.Closed())
	assert.Equal(t, true, tlsConn.AskedToBeClosed)

	//And read over already closed connection
	read, count, err = cn.ReadWithContext(buffer, maxSize, testMessageTerminator[0])
	assert.True(t, errors.Is(err, conn.ErrReaderAlreadyClosed))
	assert.Empty(t, read)
	assert.Empty(t, count)
}

func TestConnectionReadingEOF(t *testing.T) {
	tlsConn := &mock.MockTLSConnection{
		MWR: mock.MockWriteReader{
			Bytes:     []byte("Hello there, General Kenobi!"), //No terminator, so it will go until stopped
			ReturnEOF: true,
		},
	}

	cn, err := conn.NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	buffer := 50
	maxSize := 5

	read, count, err := cn.ReadWithContext(buffer, maxSize, testMessageTerminator[0])
	assert.Equal(t, true, errors.Is(err, conn.ErrStreamClosed))

	assert.Equal(t, "", string(read))
	assert.Equal(t, 0, count)

	sent, rec, errs := cn.Stats()
	assert.Equal(t, 0, rec)  //Larger buffer reads whole message before there is an error
	assert.Equal(t, 0, errs) //Got exactly 1 error
	assert.Equal(t, 0, sent)
}

func TestConnectionReadingStats(t *testing.T) {
	tlsConn := &mock.MockTLSConnection{
		MWR: mock.MockWriteReader{
			Bytes:            []byte("Hello there, General Kenobi!"), //No terminator, so it will go until stopped
			DontReturEOFEver: true,                                   //To mock endless stream with opened client
		},
	}

	cn, err := conn.NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	require.NoError(t, err)

	cn.AddRecBytes(3333)
	cn.AddErrors(15)
	cn.AddSentBytes(8374)

	sent, rec, errs := cn.Stats()
	assert.Equal(t, cn.Received(), rec)
	assert.Equal(t, cn.Errors(), errs)
	assert.Equal(t, cn.Sent(), sent)
	assert.Equal(t, cn.Received(), 3333)
	assert.Equal(t, cn.Errors(), 15)
	assert.Equal(t, cn.Sent(), 8374)

	cn.DropOldStats()

	sent, rec, errs = cn.Stats()
	assert.Equal(t, 0, rec)
	assert.Equal(t, 0, errs)
	assert.Equal(t, 0, sent)

}
