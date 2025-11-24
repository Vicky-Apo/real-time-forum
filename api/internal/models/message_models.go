package models

import "time"

// Message represents a chat message between two users
type Message struct {
	MessageID   string    `json:"message_id"`
	SenderID    string    `json:"sender_id"`
	SenderName  string    `json:"sender_name"`
	RecipientID string    `json:"recipient_id"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	IsRead      bool      `json:"is_read"`
}

// SendMessageRequest is the payload for sending a message via HTTP
type SendMessageRequest struct {
	RecipientID string `json:"recipient_id"`
	Content     string `json:"content"`
}

// SendMessageResponse is returned after successfully sending a message
type SendMessageResponse struct {
	MessageID string    `json:"message_id"`
	CreatedAt time.Time `json:"created_at"`
}

// GetMessagesResponse contains paginated message history
type GetMessagesResponse struct {
	Messages []Message `json:"messages"`
	HasMore  bool      `json:"has_more"`
}
