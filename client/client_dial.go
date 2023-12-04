package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/lazybark/go-tls-server/conn"
)

// DialTo dials to specified server and port using cert if provided.
// If cert is not provided and server has self-signed cert, DialTo will return
// 'certificate signed by unknown authority' error.
func (c *Client) DialTo(address string, port int, cert string) error {
	var config tls.Config
	if cert != "" {
		certificate, err := os.ReadFile(cert)
		if err != nil {
			return fmt.Errorf("[Client] unable to read file: %w", err)
		}
		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(certificate); !ok {
			return fmt.Errorf("[Client] unable to parse cert from %s: %w", cert, err)
		}
		config = tls.Config{RootCAs: certPool}
	}

	tlsConn, err := tls.DialWithDialer(&net.Dialer{Timeout: 3 * time.Second}, "tcp", fmt.Sprintf("%s:%d", address, port), &config)
	if err != nil {
		return fmt.Errorf("[Client] unable to dial to %s:%d: %w", address, port, err)
	}

	// We reset data in case client was used before.
	c.connCount++
	c.isClosed = false
	c.isClosedWithError = false
	c.host = address
	// Clean stats in case DropOldStats is true.
	if c.conf.DropOldStats && c.connCount > 0 {
		c.conn.DropOldStats()
	}

	cn, err := conn.NewConnection(tlsConn.RemoteAddr(), tlsConn, c.conf.MessageTerminator)
	if err != nil {
		return fmt.Errorf("[Client][Dial] error making connection for %v: %w", tlsConn.RemoteAddr(), err)
	}
	c.conn = cn

	go c.controller()
	go c.reader()

	return nil
}
