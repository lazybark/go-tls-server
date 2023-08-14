package main

import (
	"fmt"
	"log"

	"github.com/lazybark/go-tls-server/v3/server"
)

func main() {
	conf := &server.Config{KeepOldConnections: 1, HttpStatMode: true, HttpStatAddr: "localhost:8080"}
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
					err = s.SendString(conn, "Got ya!")
					if err != nil {
						log.Fatal(err)
					}
				}
			}()
		}
	}
}
