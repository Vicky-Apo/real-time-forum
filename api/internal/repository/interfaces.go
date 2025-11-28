package repository

import (
	"time"

	"real-time-forum/internal/models"
	"real-time-forum/internal/utils"
)

// UserRepositoryInterface defines the contract for user-related database operations
// This interface allows handlers to depend on behavior rather than concrete implementation
type UserRepositoryInterface interface {
	// CreateUser creates a new user in the database
	CreateUser(reg models.UserRegistration) (*models.User, error)

	// GetUserBySessionID retrieves a user by their session's user ID
	GetUserBySessionID(id string) (*models.User, error)

	// GetUserByNicknameOrEmail finds a user by nickname or email (case-insensitive)
	GetUserByNicknameOrEmail(identifier string) (*models.User, error)

	// GetAuthByUserID retrieves user authentication data (password hash)
	GetAuthByUserID(userID string) (*models.UserPassword, error)

	// Authenticate validates user credentials and returns the user if valid
	Authenticate(login models.UserLogin) (*models.User, error)

	// GetCurrentUser retrieves a user by their user ID
	GetCurrentUser(userID string) (*models.User, error)
}

// SessionRepositoryInterface defines the contract for session-related database operations
type SessionRepositoryInterface interface {
	// CreateSession creates a new session for a user
	CreateSession(userID, ipAddress string) (*models.Session, error)

	// GetBySessionID retrieves a session by its ID and validates expiration
	GetBySessionID(sessionID string) (*models.Session, error)

	// DeleteSession removes a session from the database
	DeleteSession(sessionID string) error

	// UpdateSessionIP updates the IP address for an existing session
	UpdateSessionIP(sessionID, newIP string) error
}

// PostsRepositoryInterface defines the contract for post-related database operations
type PostsRepositoryInterface interface {
	// CreatePost creates a new post with categories
	CreatePost(userID string, content string, categoryIDs []string) (*models.CreatePostResponse, error)

	// UpdatePost updates an existing post (content and categories)
	UpdatePost(postID, userID, content string, categoryIDs []string) error

	// DeletePost deletes a post (requires ownership)
	DeletePost(postID, userID string) error

	// GetPostByID retrieves a single post by its ID
	GetPostByID(postID string, userID string) (*models.Post, error)

	// GetAllPosts retrieves all posts with pagination and sorting
	GetAllPosts(limit, offset int, userID string, options utils.SortOptions) ([]*models.Post, error)

	// GetPostsByCategory retrieves posts by category with pagination and sorting
	GetPostsByCategory(categoryID string, limit, offset int, userID string, options utils.SortOptions) ([]*models.Post, error)

	// GetPostsByUser retrieves posts created by a specific user
	GetPostsByUser(targetUserID string, limit, offset int, userID *string, options utils.SortOptions) ([]*models.Post, error)

	// GetCountTotalPosts returns the total number of posts
	GetCountTotalPosts() (int, error)

	// GetCountPostByCategory returns the number of posts in a category
	GetCountPostByCategory(categoryID string) (int, error)

	// GetCountPostByUser returns the number of posts created by a user
	GetCountPostByUser(userID string) (int, error)

	// GetCountCommentedPostByUser returns the number of posts a user has commented on
	GetCountCommentedPostByUser(userID string) (int, error)

	// GetCountLikedPostByUser returns the number of posts liked by a user
	GetCountLikedPostByUser(userID string) (int, error)
}

// CategoryRepositoryInterface defines the contract for category-related database operations
type CategoryRepositoryInterface interface {
	// GetCategoryID validates that a category exists and returns its ID
	GetCategoryID(name string) (string, error)

	// GetAllCategories retrieves all categories
	GetAllCategories() ([]models.Category, error)
}

// CommentRepositoryInterface defines the contract for comment-related database operations
type CommentRepositoryInterface interface {
	// CreateComment creates a new comment on a post
	CreateComment(postID, userID, content string) (*models.CreateCommentResponse, error)

	// UpdateComment updates an existing comment (requires ownership)
	UpdateComment(commentID, userID, content string) error

	// DeleteComment deletes a comment (requires ownership)
	DeleteComment(commentID, userID string) error

	// GetCommentsByPostID retrieves comments for a specific post with pagination and sorting
	GetCommentsByPostID(postID string, limit, offset int, userID string, options utils.SortOptions) ([]*models.Comment, error)

	// GetCommentCountByPost returns the number of comments on a post
	GetCommentCountByPost(postID string) (int, error)
}

// MessageRepositoryInterface defines the contract for message-related database operations
type MessageRepositoryInterface interface {
	// SaveMessage saves a new message to the database
	SaveMessage(senderID, recipientID, content string) (*models.SendMessageResponse, error)

	// GetMessages retrieves message history between two users with pagination
	GetMessages(currentUserID, otherUserID string, limit int, beforeTimestamp *time.Time) (*models.GetMessagesResponse, error)

	// MarkMessagesAsRead marks messages from a specific sender as read
	MarkMessagesAsRead(currentUserID, senderID string) error

	// GetUnreadCount gets the count of unread messages for a user
	GetUnreadCount(userID string) (int, error)

	// GetConversations returns all conversations for a user
	GetConversations(currentUserID string) ([]models.Conversation, error)
}

// Compile-time proof that concrete types implement the interfaces
// If these assignments fail to compile, it means the concrete type doesn't implement the interface
var (
	_ UserRepositoryInterface     = (*UserRepository)(nil)
	_ SessionRepositoryInterface  = (*SessionRepository)(nil)
	_ PostsRepositoryInterface    = (*PostsRepository)(nil)
	_ CategoryRepositoryInterface = (*CategoryRepository)(nil)
	_ CommentRepositoryInterface  = (*CommentRepository)(nil)
	_ MessageRepositoryInterface  = (*MessageRepository)(nil)
)
