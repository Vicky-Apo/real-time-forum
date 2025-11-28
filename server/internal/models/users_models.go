package models

import "time"

// User represents public user information
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserAuth represents internal user data including authentication
type UserPassword struct {
	UserID       string `json:"user_id"`
	PasswordHash string `json:"-"` // Never expose in JSON
}

// UserRegistration - Registration form data 
type UserRegistration struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"` 
}

// UserLogin is used for login requests
type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
