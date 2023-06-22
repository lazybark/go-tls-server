package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"time"
)

// DialTo dials to specified server and port using cert if provided.
// If cert is not provided and server has self-signed cert, DialTo will return
// 'certificate signed by unknown authority' error
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

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 3 * time.Second}, "tcp", fmt.Sprintf("%s:%d", address, port), &config)
	if err != nil {
		return fmt.Errorf("[Client] unable to dial to %s:%d: %w", address, port, err)
	}
	c.connCount++
	c.isClosed = false
	c.isClosedWithError = false
	c.host = address
	//Clean stats in case DropOldStats is true
	if c.conf.DropOldStats && c.connCount > 0 {
		c.br = 0
		c.bs = 0
		c.errors = 0
	}

	c.conn = conn

	go c.Controller()
	go c.Reader()

	return nil
}

// Controller stops client in case stop signal recieved
func (c *Client) Controller() {
	for {
		select {
		case d := <-c.ClientDoneChan:
			if d {
				c.close(false)
				return
			}
		case <-c.ctx.Done():
			c.close(false)
			return
		}

	}
}

// Reader infinitely reads messages from opened connection
func (c *Client) Reader() {
	for {
		b, n, err := c.ReadWithContext()
		if err != nil {
			if !c.conf.SuppressErrors {
				c.ErrChan <- fmt.Errorf("[Reader] error reading from %s -> %w", c.host, err)
			}
			return
		}
		c.MessageChan <- &Message{length: n, bytes: b}
	}
}
