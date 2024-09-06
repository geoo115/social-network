package router

import (
	"net/http"
	"Social/pkg/api/handlers"
	"github.com/gorilla/websocket"
)

func HandleChatRoutes(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    upgrader := websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            return true 
        },
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
        return
    }

    // Ensure that HandleWebSocket is passed the connection
    handlers.HandleWebSocket(conn)
}