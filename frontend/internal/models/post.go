package models

import (
	"time"
)

type Post struct {
	ID         string         `json:"post_id"`
	UserID     string         `json:"user_id"`
	Username   string         `json:"username"`
	Categories []PostCategory `json:"categories"`
	Content    string         `json:"post_content"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  *time.Time     `json:"updated_at,omitempty"`

	// Aggregated metrics
	LikeCount    int `json:"like_count"`
	DislikeCount int `json:"dislike_count"`
	CommentCount int `json:"comment_count"`

	// User context
	UserReaction *int `json:"user_reaction,omitempty"` // nil, 1=like, 2=dislike
	IsOwner      bool `json:"is_owner,omitempty"`      // can current user edit/delete

	Images []PostImage `json:"images"`
}


type PostImage struct {
	ImageID          string    `json:"image_id"`
	PostID           string    `json:"post_id"`
	ImageURL         string    `json:"image_url"`
	OriginalFilename string    `json:"original_filename"`
	UploadedAt       time.Time `json:"uploaded_at"`
}
