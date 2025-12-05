package models

import "time"

// MessageImage represents an image attached to a message
type MessageImage struct {
	ImageID          string    `json:"image_id"`
	MessageID        string    `json:"message_id"`
	ImageURL         string    `json:"image_url"`
	OriginalFilename string    `json:"original_filename"`
	UploadedAt       time.Time `json:"uploaded_at"`
}

// Message represents a chat message between two users
type Message struct {
	MessageID   string         `json:"message_id"`
	SenderID    string         `json:"sender_id"`
	SenderName  string         `json:"sender_name"`
	RecipientID string         `json:"recipient_id"`
	Content     string         `json:"content"`
	CreatedAt   time.Time      `json:"created_at"`
	IsRead      bool           `json:"is_read"`
	Images      []MessageImage `json:"images"`
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

// Conversation represents a chat conversation with another user
type Conversation struct {
	UserID      string       `json:"user_id"`
	Username    string       `json:"username"`
	IsOnline    bool         `json:"is_online"`
	LastMessage *LastMessage `json:"last_message"`
	UnreadCount int          `json:"unread_count"`
}

// LastMessage represents the most recent message in a conversation
type LastMessage struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsFromMe  bool      `json:"is_from_me"`
}
