package main

import (
	"fmt"
	"log"

	"github.com/lazybark/go-tls-server/v2/server"
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
		case m := <-s.MessageChan:
			fmt.Println("Got message:", string(m.Bytes()))
			_, err = m.Conn().SendString("Got ya!")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
