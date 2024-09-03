package router

import (
	"Social/pkg/api/handlers"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleChatRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// List of allowed origins for WebSocket connections
	allowedOrigins := map[string]bool{
		"http://localhost:3000": true, // React frontend
		"http://localhost:5173": true, // Vite frontend
		// Add other allowed origins as needed
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if allowedOrigins[origin] {
				return true
			}
			log.Printf("WebSocket connection attempted from disallowed origin: %s", origin)
			return false
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	// Pass the WebSocket connection to the handler
	handlers.HandleWebSocket(conn)
}
