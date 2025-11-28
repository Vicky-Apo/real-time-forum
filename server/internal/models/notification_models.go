package models

import "time"

// Notification represents a user notification
type Notification struct {
	NotificationID     string    `json:"notification_id"`
	UserID             string    `json:"user_id"`              // who gets the notification
	TriggerUsername    string    `json:"trigger_username"`     // who caused it (e.g., "John")
	PostContentPreview string    `json:"post_content_preview"` // first 50 chars of post content
	PostID             string    `json:"post_id"`              // link to the post
	Action             string    `json:"action"`               // "liked" or "commented on"
	IsRead             bool      `json:"is_read"`
	CreatedAt          time.Time `json:"created_at"`
}

// NotificationResponse represents the API response for notifications
type NotificationResponse struct {
	Notifications []*Notification `json:"notifications"`
	TotalCount    int             `json:"total_count,omitempty"`
}
