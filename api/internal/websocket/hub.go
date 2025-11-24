package websocket

import (
	"encoding/json"
	"log"

	"real-time-forum/internal/models"
)

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	// Registered clients mapped by user ID
	Clients map[string]*Client

	// Register requests from the clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			// Register the client
			h.Clients[client.UserID] = client
			log.Printf("User %s (%s) connected. Total connections: %d", client.Nickname, client.UserID, len(h.Clients))

			// Notify all other users that this user is online
			h.BroadcastUserStatus(client.UserID, client.Nickname, "online")

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
				log.Printf("User %s (%s) disconnected. Total connections: %d", client.Nickname, client.UserID, len(h.Clients))

				// Notify all other users that this user is offline
				h.BroadcastUserStatus(client.UserID, client.Nickname, "offline")
			}
		}
	}
}

// HandleMessage processes incoming messages from clients based on event type
func (h *Hub) HandleMessage(sender *Client, msg models.WebSocketMessage) {
	switch msg.Event {
	case models.EventTypeTypingStart:
		h.handleTypingIndicator(sender, msg.Payload, true)
	case models.EventTypeTypingStop:
		h.handleTypingIndicator(sender, msg.Payload, false)
	default:
		log.Printf("Unsupported WebSocket event type from %s: %s", sender.Nickname, msg.Event)
		sender.SendError("Unsupported event type. Use HTTP POST /api/messages/send to send messages.")
	}
}

// BroadcastUserStatus broadcasts a user's online/offline status to all connected users
func (h *Hub) BroadcastUserStatus(userID, nickname, status string) {
	payload := models.UserStatusPayload{
		UserID:   userID,
		Nickname: nickname,
		Status:   status,
	}

	var event string
	if status == "online" {
		event = models.EventTypeUserOnline
	} else {
		event = models.EventTypeUserOffline
	}

	// Send to all clients except the user themselves
	for _, client := range h.Clients {
		if client.UserID != userID {
			client.SendMessage(event, payload)
		}
	}
}

// GetOnlineUsers returns a list of currently online users
func (h *Hub) GetOnlineUsers() []models.UserStatusPayload {
	users := make([]models.UserStatusPayload, 0, len(h.Clients))
	for _, client := range h.Clients {
		users = append(users, models.UserStatusPayload{
			UserID:   client.UserID,
			Nickname: client.Nickname,
			Status:   "online",
		})
	}
	return users
}

// SendMessageToUser sends a message to a specific user if they are online
// Returns true if the user is online and message was sent, false otherwise
func (h *Hub) SendMessageToUser(userID string, event string, payload interface{}) bool {
	client, ok := h.Clients[userID]
	if !ok {
		return false // User is not online
	}

	client.SendMessage(event, payload)
	return true
}

// handleTypingIndicator handles typing start/stop events
func (h *Hub) handleTypingIndicator(sender *Client, payload interface{}, isTyping bool) {
	// Parse the payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		sender.SendError("Invalid typing indicator payload")
		return
	}

	var typingPayload models.TypingIndicatorPayload
	err = json.Unmarshal(payloadBytes, &typingPayload)
	if err != nil {
		sender.SendError("Invalid typing indicator format")
		return
	}

	// Validate recipient ID
	if typingPayload.RecipientID == "" {
		sender.SendError("Recipient ID is required")
		return
	}

	// Don't allow typing indicator to yourself
	if typingPayload.RecipientID == sender.UserID {
		return
	}

	// Check if recipient is online
	recipient, ok := h.Clients[typingPayload.RecipientID]
	if !ok {
		// Recipient is not online, silently ignore (no error needed)
		return
	}

	// Send typing notification to recipient
	notification := models.TypingNotificationPayload{
		UserID:   sender.UserID,
		Nickname: sender.Nickname,
		IsTyping: isTyping,
	}

	var eventType string
	if isTyping {
		eventType = models.EventTypeTypingStart
	} else {
		eventType = models.EventTypeTypingStop
	}

	recipient.SendMessage(eventType, notification)

	log.Printf("Typing indicator: %s → %s (is_typing: %v)", sender.Nickname, recipient.Nickname, isTyping)
}
