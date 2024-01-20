package main

import (
	"log"

	"github.com/lazybark/go-tls-server/server"
)

func main() {
	conf := &server.Config{
		SuppressErrors:          false,
		MaxMessageSize:          0,
		MessageTerminator:       '\n',
		BufferSize:              128, //nolint:gomnd // It's OK
		KeepOldConnections:      1,
		KeepInactiveConnections: 4320, //nolint:gomnd // It's OK
		ErrorPrefix:             "MY_SERVER",
	}

	tlsServer, err := server.New("localhost", `certs/cert.pem`, `certs/key.pem`, conf)
	if err != nil {
		log.Fatal(err)
	}

	go tlsServer.Listen("5555")

	for {
		select {
		case err, ok := <-tlsServer.ErrChan():
			if !ok {
				return
			}

			log.Println(err)

		case conn, ok := <-tlsServer.ConnChan():
			if !ok {
				return
			}

			log.Println(conn.Address())

			go func() {
				for m := range conn.MessageChanRead() {
					log.Println("Got message:", string(m.Bytes()))

					err = tlsServer.SendString(conn, "Got ya!")
					if err != nil {
						log.Fatal(err)
					}
				}
			}()
		}
	}
}
