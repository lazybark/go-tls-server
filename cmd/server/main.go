package main

import (
	"context"
	"errors"
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

	tlsServer, err := server.New(context.Background(), "localhost", `certs/cert.pem`, `certs/key.pem`, conf)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err = tlsServer.Error()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = tlsServer.Listen("5555")
	if err != nil {
		log.Fatal(err)
	}

	for tlsServer.Next() {
		connection, err := tlsServer.AcceptConnection()
		if err != nil {
			if !errors.Is(err, server.ErrServerClosed) {
				log.Fatal(err)
			}
		}

		go func() {
			for connection.Next() {
				message, err := connection.GetMessage()
				if err != nil {
					return
				}

				log.Println("Got message:", string(message.Bytes()))

				err = tlsServer.SendString(connection, "Got ya!")
				if err != nil {
					log.Fatal(err)
				}
			}
		}()

	}
}
