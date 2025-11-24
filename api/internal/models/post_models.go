package models

import "time"

type Post struct {
	ID         string         `json:"post_id"`
	UserID     string         `json:"user_id"`
	Nickname   string         `json:"nickname"`
	Categories []PostCategory `json:"categories"`
	Content    string         `json:"post_content"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  *time.Time     `json:"updated_at,omitempty"`

	// Aggregated metrics
	CommentCount int `json:"comment_count"`
}

// CreatePostRequest - Post creation payload
type CreatePostRequest struct {
	CategoryNames []string `json:"category_names" binding:"required,min=1,max=5"`
	Content       string   `json:"content" binding:"required,min=10,max=5000"`
}
