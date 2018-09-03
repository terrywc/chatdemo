package main

import (
	"fmt"
	"net"
)

type client struct {
	socket net.Conn
	data   chan []byte
}

func (client *client) receive() {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED: " + string(message))
		}
	}
}
