package models

import "time"

// User represents public user information
type User struct {
	ID        string    `json:"id"`
	Nickname  string    `json:"nickname"`
	Age       int       `json:"age"`
	Gender    string    `json:"gender"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserAuth represents internal user data including authentication
type UserPassword struct {
	UserID       string `json:"user_id"`
	PasswordHash string `json:"-"` // Never expose in JSON
}

// UserRegistration - Registration form data (updated with new required fields)
type UserRegistration struct {
	Nickname        string `json:"nickname"`
	Age             int    `json:"age"`
	Gender          string `json:"gender"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// UserLogin is used for login requests (accepts nickname or email)
type UserLogin struct {
	Identifier string `json:"identifier" binding:"required"` // Can be nickname or email
	Password   string `json:"password" binding:"required"`
}
