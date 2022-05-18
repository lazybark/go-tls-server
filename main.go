package main

import (
	"fmt"
	"log"

	v1 "github.com/lazybark/go-tls-server/v1"
)

func main() {
	conf := &v1.Config{KeepOldConnections: 1}
	s, err := v1.New("localhost", `C:\Users\serge\Desktop\git repos\lazybark\go-tls-server\certs\cert.pem`, `C:\Users\serge\Desktop\git repos\lazybark\go-tls-server\certs\key.pem`, conf)
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
		}
	}
}
