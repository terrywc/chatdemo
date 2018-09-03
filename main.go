package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func startServerMode() {
	fmt.Println("Starting server...")
	listener, error := net.Listen("tcp", "localhost:12345")
	if error != nil {
		fmt.Println(error)
	}
	manager := clientmanager{
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
	go manager.start()
	for {
		connection, _ := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}
		client := &client{socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client)
		go manager.send(client)
	}
}

func startClientMode() {
	fmt.Println("Starting client...")
	connection, error := net.Dial("tcp", "localhost:12345")
	if error != nil {
		fmt.Println(error)
	}
	client := &client{socket: connection}
	go client.receive()
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		connection.Write([]byte(strings.TrimRight(message, "\n")))
	}
}

func main() {
	flagMode := flag.String("mode", "server", "start in client or server mode")
	flag.Parse()
	if strings.ToLower(*flagMode) == "server" {
		startServerMode()
	} else {
		startClientMode()
	}
}
