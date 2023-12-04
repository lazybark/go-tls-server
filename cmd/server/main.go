package main

import (
	"fmt"
	"log"

	"github.com/lazybark/go-tls-server/server"
)

func main() {
	conf := &server.Config{KeepOldConnections: 1}
	s, err := server.New("localhost", `certs/cert.pem`, `certs/key.pem`, conf)
	if err != nil {
		log.Fatal(err)
	}

	go s.Listen("5555")

	for {
		select {
		case err, ok := <-s.ErrChan():
			if !ok {
				return
			}

			fmt.Println(err)

		case conn, ok := <-s.ConnChan():
			if !ok {
				return
			}

			fmt.Println(conn.Address())

			go func() {
				for m := range conn.MessageChanRead() {
					fmt.Println("Got message:", string(m.Bytes()))
					err = s.SendString(conn, "Got ya!")
					if err != nil {
						log.Fatal(err)
					}
				}
			}()
		}
	}
}
