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
	EventTypeTypingStart    = "typing_start"
	EventTypeTypingStop     = "typing_stop"
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
	SenderID   string         `json:"sender_id"`   // User ID of the sender
	SenderName string         `json:"sender_name"` // Username of the sender
	Content    string         `json:"content"`     // Message content
	SentAt     time.Time      `json:"sent_at"`     // Timestamp when message was sent
	Images     []MessageImage `json:"images"`      // Attached images
}

// ErrorPayload represents an error message
type ErrorPayload struct {
	Message string `json:"message"` // Error message
}

// UserStatusPayload represents a user's online/offline status
type UserStatusPayload struct {
	UserID   string `json:"user_id"`   // User ID
	Username string `json:"username"`  // Username
	Status   string `json:"status"`    // "online" or "offline"
}

// TypingIndicatorPayload represents a typing indicator event (sent from client)
type TypingIndicatorPayload struct {
	RecipientID string `json:"recipient_id"` // User ID of who should see the typing indicator
}

// TypingNotificationPayload represents a typing notification (sent to recipient)
type TypingNotificationPayload struct {
	UserID   string `json:"user_id"`   // User ID of who is typing
	Username string `json:"username"`  // Username of who is typing
	IsTyping bool   `json:"is_typing"` // true = started typing, false = stopped typing
}
