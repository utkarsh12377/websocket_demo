package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {

	serverURL := "ws://localhost:8080/ws"

	fmt.Println("Connecting to chat server...")

	conn, _, err := websocket.DefaultDialer.Dial(
		serverURL,
		nil,
	)

	if err != nil {
		log.Fatal(
			"Connection failed:",
			err,
		)
	}

	defer conn.Close()

	fmt.Println(
		"Connected to chat server!",
	)

	reader := bufio.NewReader(
		os.Stdin,
	)

	
	fmt.Print("Enter your name: ")

	name, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	name = strings.TrimSpace(name)

	
	err = conn.WriteMessage(
		websocket.TextMessage,
		[]byte(name),
	)

	if err != nil {
		log.Fatal(
			"Failed to send username:",
			err,
		)
	}

	fmt.Println()
	fmt.Println(
		"Welcome,",
		name,
	)

	fmt.Println(
		"Type messages below.",
	)

	fmt.Println(
		"Type 'exit' to leave.",
	)

	fmt.Println()

	
	go func() {

		for {

			_, message, err :=
				conn.ReadMessage()

			if err != nil {

				fmt.Println(
					"\nDisconnected from server.",
				)

				return
			}

			fmt.Println(
				"\n" + string(message),
			)

			fmt.Print("You: ")
		}

	}()

	
	for {

		fmt.Print("You: ")

		message, err :=
			reader.ReadString('\n')

		if err != nil {
			return
		}

		message =
			strings.TrimSpace(message)

		if message == "" {
			continue
		}

		if message == "exit" {

			fmt.Println(
				"Leaving chat...",
			)

			return
		}

		err = conn.WriteMessage(
			websocket.TextMessage,
			[]byte(message),
		)

		if err != nil {

			fmt.Println(
				"Send error:",
				err,
			)

			return
		}
	}
}