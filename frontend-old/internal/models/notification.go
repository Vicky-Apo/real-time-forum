package models

import (
	"time"
)

// Notification represents a user notification (matches backend exactly)
type Notification struct {
	NotificationID     string    `json:"notification_id"`
	UserID             string    `json:"user_id"`              // who gets the notification
	TriggerUsername    string    `json:"trigger_username"`     // who caused it (e.g., "John")
	PostContentPreview string    `json:"post_content_preview"` // first 50 chars of post content
	PostID             string    `json:"post_id"`              // link to the post
	Action             string    `json:"action"`               // "liked" or "commented on" or "disliked"
	IsRead             bool      `json:"is_read"`
	CreatedAt          time.Time `json:"created_at"`
}

// NotificationResponse represents the API response for notifications (matches backend exactly)
type NotificationResponse struct {
	Notifications []*Notification `json:"notifications"`
	TotalCount    int             `json:"total_count,omitempty"`
}

// NotificationPageData represents data for the notifications template page
type NotificationPageData struct {
	Notifications []*Notification `json:"notifications"`
	User          *User           `json:"user,omitempty"` // Current logged-in user
	UnreadCount   int             `json:"unread_count"`
	TotalCount    int             `json:"total_count"`
}

// Helper method to check if notification is unread
func (n *Notification) IsUnread() bool {
	return !n.IsRead
}
