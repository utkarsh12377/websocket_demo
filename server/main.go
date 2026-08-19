package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	name string
}

type Hub struct {
	clients map[*Client]bool
	mutex   sync.Mutex
}

var hub = Hub{
	clients: make(map[*Client]bool),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	defer conn.Close()

	fmt.Println("New client connected")

	_, username, err := conn.ReadMessage()

	if err != nil {
		fmt.Println("Error reading username:", err)
		return
	}

	client := &Client{
		conn: conn,
		name: string(username),
	}

	hub.mutex.Lock()
	hub.clients[client] = true
	hub.mutex.Unlock()

	fmt.Println(client.name, "joined the chat")

	broadcast(
		fmt.Sprintf(
			"%s joined the chat",
			client.name,
		),
	)

	for {

		_, message, err := conn.ReadMessage()

		if err != nil {

			fmt.Println(
				client.name,
				"disconnected",
			)

			hub.mutex.Lock()
			delete(hub.clients, client)
			hub.mutex.Unlock()

			broadcast(
				fmt.Sprintf(
					"%s left the chat",
					client.name,
				),
			)

			break
		}

		formattedMessage := fmt.Sprintf(
			"%s: %s",
			client.name,
			string(message),
		)

		fmt.Println(formattedMessage)

		broadcast(formattedMessage)
	}
}

func broadcast(message string) {

	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	for client := range hub.clients {

		err := client.conn.WriteMessage(
			websocket.TextMessage,
			[]byte(message),
		)

		if err != nil {
			fmt.Println(
				"Broadcast error:",
				err,
			)
		}
	}
}

func main() {

	http.HandleFunc(
		"/ws",
		websocketHandler,
	)

	fmt.Println(
		"Chat server running on ws://localhost:8080/ws",
	)

	err := http.ListenAndServe(
		":8080",
		nil,
	)

	if err != nil {
		fmt.Println(
			"Server error:",
			err,
		)
	}
}
