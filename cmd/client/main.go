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

	done := make(chan bool)

	go func() {
		for err := range c.ErrChan {
			fmt.Println(err)
		}
	}()

	go func() {
		for m := range c.MessageChan {
			fmt.Println("Got message:", string(m.Bytes()))
		}
		done <- true
	}()

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

	<-done
}
