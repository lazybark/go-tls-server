package conn

import (
	"fmt"
	"testing"

	"github.com/lazybark/go-helpers/gen"
	"github.com/lazybark/go-helpers/mock"
)

var (
	readWithContextBenchmarkResult []byte
	sendByteBenchmarkResult        int
	sendStringBenchmarkResult      int
)

func FailTest(b *testing.B, err error) {
	if err != nil {
		b.Error(err)
	}
}

// Small buffer affects reading speed
func BenchmarkConnectionCorrectReading(b *testing.B) {
	tlsConn := &mock.MockTLSConnection{
		MWR: mock.MockWriteReader{},
	}

	cn, _ := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')

	bLens := []int{5, 128, 1024}
	maxSize := 5120

	for _, bl := range bLens {
		for n := 1024; n <= 5120; n += 1024 {
			b.Run(fmt.Sprintf("size_%d_buffer_%d", n, bl), func(b *testing.B) {
				tlsConn.MWR.Bytes = []byte(gen.GenerateRandomString(n))
				for i := 0; i < b.N; i++ {
					readWithContextBenchmarkResult, _, _ = cn.ReadWithContext(bl, maxSize, testMessageTerminator[0])

					tlsConn.MWR.SetLastRead(0)
				}
			})
		}
	}
}

func BenchmarkConnectionCorrectByteSending(b *testing.B) {
	var send []string
	var s string
	tlsConn := &mock.MockTLSConnection{}

	for i := 256; i <= 1024; i += 256 {
		s = gen.GenerateRandomString(i)

		send = append(send, s)
	}

	cn, _ := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	for _, str := range send {
		//Sending bytes, so counting also bytes, not chars
		b.Run(fmt.Sprintf("input_size_%d", len(str)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sendByteBenchmarkResult, _ = cn.SendByte([]byte(str))
			}
		})
	}
}

func BenchmarkConnectionCorrectStringSending(b *testing.B) {
	var send []string
	var s string
	tlsConn := &mock.MockTLSConnection{}

	for i := 256; i <= 1024; i += 256 {
		s = gen.GenerateRandomString(i)

		send = append(send, s)
	}

	cn, _ := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	for _, str := range send {
		//Sending bytes, so counting also bytes, not chars
		b.Run(fmt.Sprintf("input_size_%d", len(str)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sendStringBenchmarkResult, _ = cn.SendString(str)
			}
		})
	}
}

func BenchmarkConnectionStatsLastAct(b *testing.B) {
	tlsConn := &mock.MockTLSConnection{}

	cn, _ := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')

	for i := 0; i < b.N; i++ {
		cn.setLastAct()
	}
}

func BenchmarkConnectionStatsAddRecBytes(b *testing.B) {
	tlsConn := &mock.MockTLSConnection{}

	cn, _ := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	for n := 256; n <= 1024; n += 256 {
		b.Run(fmt.Sprintf("size_%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cn.addRecBytes(n)
			}
		})
	}
}

func BenchmarkConnectionStatsAddSentBytes(b *testing.B) {
	tlsConn := &mock.MockTLSConnection{}

	cn, _ := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	for n := 256; n <= 1024; n += 256 {
		b.Run(fmt.Sprintf("size_%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cn.addSentBytes(n)
			}
		})
	}
}

func BenchmarkConnectionStatsAddErrors(b *testing.B) {
	tlsConn := &mock.MockTLSConnection{}

	cn, _ := NewConnection(tlsConn.RemoteAddr(), tlsConn, '\n')
	for n := 256; n <= 1024; n += 256 {
		b.Run(fmt.Sprintf("size_%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cn.addErrors(n)
			}
		})
	}
}
