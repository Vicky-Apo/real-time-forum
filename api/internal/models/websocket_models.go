package models

import "time"

// WebSocket Event Types
const (
	// DEPRECATED: EventTypeSendMessage - No longer supported. Use HTTP POST /api/messages/send instead.
	EventTypeSendMessage    = "send_message"
	EventTypeReceiveMessage = "receive_message"
	EventTypeError          = "error"
	EventTypeUserOnline     = "user_online"
	EventTypeUserOffline    = "user_offline"
)

// WebSocketMessage represents the generic WebSocket message structure
// Note: WebSocket is currently receive-only for messages. Use HTTP API to send messages.
type WebSocketMessage struct {
	Event   string      `json:"event"`   // Event type (receive_message, user_online, user_offline, error)
	Payload interface{} `json:"payload"` // Event-specific payload
}

// DEPRECATED: SendMessagePayload - No longer supported via WebSocket.
// Use HTTP POST /api/messages/send with SendMessageRequest instead.
type SendMessagePayload struct {
	RecipientID string `json:"recipient_id"` // User ID of the recipient
	Content     string `json:"content"`      // Message content
}

// ReceiveMessagePayload represents the payload for receiving a message
type ReceiveMessagePayload struct {
	SenderID   string    `json:"sender_id"`   // User ID of the sender
	SenderName string    `json:"sender_name"` // Nickname of the sender
	Content    string    `json:"content"`     // Message content
	SentAt     time.Time `json:"sent_at"`     // Timestamp when message was sent
}

// ErrorPayload represents an error message
type ErrorPayload struct {
	Message string `json:"message"` // Error message
}

// UserStatusPayload represents a user's online/offline status
type UserStatusPayload struct {
	UserID   string `json:"user_id"`   // User ID
	Nickname string `json:"nickname"`  // User nickname
	Status   string `json:"status"`    // "online" or "offline"
}
