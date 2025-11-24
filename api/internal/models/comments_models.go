package models

import "time"

type Comment struct {
	ID        string     `json:"comment_id"`
	UserID    string     `json:"user_id"`
	Nickname  string     `json:"nickname"`
	PostID    string     `json:"post_id"`
	Content   string     `json:"comment_content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// CreateCommentRequest is used for creating a new comment
type CreateCommentRequest struct {
	PostID  string `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required,min=10,max=500"`
}
