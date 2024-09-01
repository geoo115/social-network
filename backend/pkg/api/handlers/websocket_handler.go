package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan models.Chat)       // Broadcast channel
var mutex = &sync.Mutex{}

func HandleWebSocket(conn *websocket.Conn) {
	// Register the new client
	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	// Ensure connection is closed when the function ends
	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
		conn.Close()
	}()

	// Listen for incoming messages
	for {
		var message models.Chat
		err := conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Printf("Unexpected WebSocket closure: %v", err)
				break
			}
			log.Printf("Error reading JSON: %v", err)
			continue
		}

		// Log the received message
		log.Printf("Received WebSocket message: %+v", message)

		// Ensure that SenderID and RecipientID are populated correctly
		if message.CreatedAt.IsZero() {
			message.CreatedAt = time.Now()
		}

		// Store the message in the database
		if err := services.SendMessage(message); err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}

		// Broadcast the message to all connected clients
		broadcast <- message
	}
}

// HandleMessages listens for messages on the broadcast channel and sends them to all clients
func HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		message := <-broadcast

		// Send it out to every client that is currently connected
		mutex.Lock()
		for client := range clients {
			if err := client.WriteJSON(message); err != nil {
				log.Printf("Error writing JSON: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
