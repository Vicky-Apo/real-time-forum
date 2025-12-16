package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"real-time-forum/internal/middleware"
	"real-time-forum/internal/models"
	ws "real-time-forum/internal/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow connections from any origin (adjust for production)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user from context
		user := middleware.GetCurrentUser(r)

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		// Create a new client
		client := &ws.Client{
			Hub:      hub,
			Conn:     conn,
			UserID:   user.ID,
			Username: user.Username,
			Send:     make(chan models.WebSocketMessage, 256),
		}

		// Register the client with the hub
		hub.Register <- client

		// Start goroutines for reading and writing
		go client.WritePump()
		go client.ReadPump()
	}
}
