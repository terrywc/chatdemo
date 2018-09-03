package main

import (
	"fmt"
)

type clientmanager struct {
	clients    map[*client]bool
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
}

func (manager *clientmanager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Println("Added new connection!")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				fmt.Println("A connection has terminated!")
			}
		case message := <-manager.broadcast:
			for connection := range manager.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

func (manager *clientmanager) receive(client *client) {
	for {
		message := make([]byte, 1024)
		length, err := client.socket.Read(message)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED: " + string(message) + " Client: " + string(client.socket.RemoteAddr().String()))
			manager.broadcast <- message
		}
	}
}

func (manager *clientmanager) send(client *client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			client.socket.Write(message)
		}
	}
}
