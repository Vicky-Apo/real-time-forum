package models

// UserRegistration - Registration form data (matches backend exactly)
type UserRegistration struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"` // Added confirm password field
}

// UserLogin - Login form data (matches backend exactly)
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
