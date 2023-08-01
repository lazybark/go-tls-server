# go-tls-server
[![Test](https://github.com/lazybark/go-tls-server/actions/workflows/test.yml/badge.svg)](https://github.com/lazybark/go-tls-server/actions/workflows/test.yml)

go-tls-server is a small lib to create client-server apps using tls.Conn. It uses standard libs to create stream-like message exchange protected by TLS. Every message ends with a terminator (:robot:) symbol and the main idea is to read from connection until :robot: appears, then process what we read and repeat reading. This way we can create apps that control their behaviour using any possible custom protocol/message set.

A practical example of how it works you can find in [go-cloud-sync](https://github.com/lazybark/go-cloud-sync).

Cert & key for server & client can be generated via [go-cert-generator](https://github.com/lazybark/go-cert-generator).

Server parameters:
* SuppressErrors (bool) - prevents server from sending errors into ErrChan
* MaxMessageSize (int) - sets max length of one message in bytes
* MessageTerminator (byte) - sets byte value that marks message end of the message in stream
* BufferSize (int) - regulates buffer length to read incoming message
* KeepOldConnections (int) - prevents server from dropping closed connection for N minutes after it has been closed
* KeepInactiveConnections (int) - makes server close connection that had no activity for N mins

Client parameters:
* SuppressErrors (bool) - prevents client from sending errors into ErrChan
* MaxMessageSize (int) - sets max length of one message in bytes
* MessageTerminator (byte) - sets byte value that marks message end of the message in stream
* BufferSize (int) - regulates buffer length to read incoming message
* DropOldStats (bool) - make client to set all sent/recieved bytes & errors to zero before opening new connection

### Control connections
Server manages connections by deleting old & inactive from connPool. So when you use similar connection pool in your project (to store client-related data), you might need to check if the connection is still active. Server stores pointers and deletes them after some period of time, but if your app stores pointers to server connections, then you will not notice the fact that connection was removed from server. It will still be accessible and if it has been closed, you will encounter an error when trying write/read. The best way to check if connection is still usable is to call Connection.Closed().

### Statistic
Both client and server have stats than can be useful. 

Server has:
* Stats(year int, month int, day int) - will return number of bytes sent/received + number of errors or an ErrNoStatForTheDay 
* StatsConnections() - will simply return current number of connections in pool
* ActiveConnetions() - total number of currently active (usable) connections
* Online() - how long the server is online

Client has:
* Stats() - will return number of bytes sent/received + number of errors

## Basic usage
Basic usage is to use server & client behind an interface or as part of bigger struct. Both return new connections and messages via channels to external calling code which means you can create routines to process new connections and messages in them (as server) or to create separate connections and communicate with server (as client).

So you just run a routine that awaits in connection channel and does some magic when new connection appears. Best way here is to add connection to your internal pool (if you need to manage it with some extra data) and then run goroutine that awaits & processes messages via connection message channel.

### Simple server code

```
package main

import (
	"fmt"
	"log"

	"github.com/lazybark/go-tls-server/v3/server"
)

func main() {
	conf := &server.Config{KeepOldConnections: 1, NotifyAboutNewConnections: true}
	s, err := server.New("localhost", `certs/cert.pem`, `certs/key.pem`, conf)
	if err != nil {
		log.Fatal(err)
	}
	go s.Listen("5555")
	for {
		select {
		case err := <-s.ErrChan:
			fmt.Println(err)
		case conn := <-s.ConnChan:
			fmt.Println(conn.Address())
			go func() {
				for m := range conn.MessageChan {
					fmt.Println("Got message:", string(m.Bytes()))
					_, err = conn.SendString("Got ya!")
					if err != nil {
						log.Fatal(err)
					}
				}
			}()
		}
	}
}

```

### Simple client code
```
package main

import (
	"fmt"
	"log"

	"github.com/lazybark/go-tls-server/v3/client"
)

func main() {
	ipsum := `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Viverra nibh cras pulvinar mattis nunc sed. Congue nisi vitae suscipit tellus. Enim sit amet venenatis urna cursus. Egestas integer eget aliquet nibh. Orci phasellus egestas tellus rutrum tellus pellentesque eu tincidunt. Feugiat vivamus at augue eget arcu dictum varius. Tincidunt praesent semper feugiat nibh sed pulvinar proin gravida. Neque gravida in fermentum et sollicitudin. Purus in massa tempor nec feugiat nisl. Vitae purus faucibus ornare suspendisse. Viverra tellus in hac habitasse. Aliquam sem et tortor consequat id porta nibh. Ipsum suspendisse ultrices gravida dictum fusce. Fermentum iaculis eu non diam phasellus. Ultrices eros in cursus turpis massa. Ut ornare lectus sit amet est placerat in. Id ornare arcu odio ut sem nulla pharetra.`
	conf := client.Config{SuppressErrors: false, MessageTerminator: '\n'}
	c := client.New(&conf)

	err := c.DialTo("localhost", 5555, `certs/cert.pem`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.SendString(ipsum)
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.SendByte([]byte{'H', 'i', '!'})
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-c.ErrChan:
			fmt.Println(err)
		case m := <-c.MessageChan:
			fmt.Println("Got message:", string(m.Bytes()))
		}
	}
}
```

## Complex usage

### Server
As an example of using the server in bigger project, you can create struct that holds config for server:
```
package server

import "github.com/lazybark/go-tls-server/v3/conn"

type LinkServer struct {
	extConnChan chan (*conn.Connection)
	extErrChan  chan (error)
	certPath    string
	keyPath     string
}

func NewServer(certPath, keyPath string) *LinkServer {
	s := &LinkServer{
		certPath: certPath,
		keyPath:  keyPath,
	}
	return s
}
```
Then create methods to init server and start:
```
package server

import (
	"fmt"

	"github.com/lazybark/go-tls-server/v3/conn"
	gts "github.com/lazybark/go-tls-server/v3/server"
)

// Init prepares server to accept connections
func (s *LinkServer) Init(extConnChan chan (*conn.Connection), extErrChan chan (error)) error {
	s.extConnChan = extConnChan
	s.extErrChan = extErrChan
	return nil
}

// Listen starts net listener
func (s *LinkServer) Listen(addr, port string) error {
	conf := &gts.Config{KeepOldConnections: 1, NotifyAboutNewConnections: true}
	srv, err := gts.New(addr, s.certPath, s.keyPath, conf)
	if err != nil {
		return fmt.Errorf("[Link][Listen]%w", err)
	}

	go srv.Listen(port)

	go func() {

		for {
			select {
			case err, ok := <-srv.ErrChan:
				if !ok {
					return
				}

				s.extErrChan <- fmt.Errorf("[Link][Listen]%w", err)

			case c, ok := <-srv.ConnChan:
				if !ok {
					return
				}

				s.extConnChan <- c
			}
		}
	}()

	return nil
}
```
Now your server can already accept connections. LinkServer struct can be then put into bigger struct that controls your app.

Same goes for client. You create some struct with parameters and client field in it.
```
package client

import (
	"github.com/lazybark/go-tls-server/v3/client"
)

// LinkClient  works with lazybark/go-tls-server to implement ISyncLinkClientV1 interface
type LinkClient struct {
	certPath   string
	akey       string
	cid        string
	serverAddr string
	serverPort int
	login      string
	pwd        string
	c          *client.Client
}

// NewClient returns new LinkClient
func NewClient(certPath string) (*LinkClient, error) {
	c := &LinkClient{certPath: certPath}

	return c, nil
}
```
Then you add config functions:
```
package client

import (
	"fmt"

	tls "github.com/lazybark/go-tls-server/v3/client"
)

// setAuth sets existing key & session ID to the connection
func (sc *LinkClient) setAuth(akey string, cid string) {
	sc.akey = akey
	sc.cid = cid
}

// Init prepares client to connect
func (sc *LinkClient) Init(port int, addr, login, pwd string) error {

	conf := tls.Config{SuppressErrors: false, MessageTerminator: '\n'}
	c := tls.New(&conf)
	sc.c = c
	sc.serverAddr = addr
	sc.serverPort = port
	sc.login = login
	sc.pwd = pwd

	err := sc.ConnectAndAuth()
	if err != nil {
		return fmt.Errorf("[SyncClient][Init]%w", err)
	}

	return nil
}

// Close closes the client
func (sc *LinkClient) Close() error {
	return sc.c.Close()
}

```
And use ending struct in your client app. For example: to wait and process one single response from server, you can create such method:
```
package client

import (
	"encoding/json"
	"fmt"
)

// Await awaits for exactly one message in the stream and returns its content
func (sc *LinkClient) Await() (MessageModel, error) {
	var m MessageModel
	ans := <-sc.c.MessageChan
	err := json.Unmarshal(ans.Bytes(), &m)
	if err != nil {
		return m, fmt.Errorf("[AwaitAnswer]%w", err)
	}
	return m, nil
}
```
It will return MessageModel, then you can process it and run another Await() after sending new data for server.


