package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func main() {

	serverURL := "ws://localhost:8080/ws"

	fmt.Println("Connecting to server...")

	conn, _, err := websocket.DefaultDialer.Dial(
		serverURL,
		nil,
	)

	if err != nil {
		log.Fatal("Connection failed:", err)
	}

	defer conn.Close()

	fmt.Println("Connected to WebSocket server!")

	message := "Hello from Go Client"

	err = conn.WriteMessage(
		websocket.TextMessage,
		[]byte(message),
	)

	if err != nil {
		log.Fatal("Write error:", err)
	}

	fmt.Println("Message sent:", message)

	messageType, response, err :=
		conn.ReadMessage()

	if err != nil {
		log.Fatal("Read error:", err)
	}

	fmt.Println(
		"Server response:",
		string(response),
	)

	fmt.Println(
		"Message type:",
		messageType,
	)
}