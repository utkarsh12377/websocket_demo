package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Incoming WebSocket connection...")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	defer conn.Close()

	fmt.Println("Client connected!")

	for {

		messageType, message, err := conn.ReadMessage()

		if err != nil {
			fmt.Println("Client disconnected")
			break
		}

		fmt.Println(
			"Received from client:",
			string(message),
		)

		response := fmt.Sprintf(
			"Server received: %s",
			string(message),
		)

		err = conn.WriteMessage(
			messageType,
			[]byte(response),
		)

		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}

func main() {

	http.HandleFunc(
		"/ws",
		websocketHandler,
	)

	fmt.Println(
		"WebSocket server running on ws://localhost:8080/ws",
	)

	err := http.ListenAndServe(
		":8080",
		nil,
	)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
