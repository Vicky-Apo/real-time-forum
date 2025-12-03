// internal/utils/session_utils.go
package utils

import (
	"net/http"

	"frontend-service/internal/models"
)

// GetUserFromSession gets user from session cookie and validates with backend
// This replaces the session.GetUserFromSession function
func GetUserFromSession(r *http.Request, authService AuthServiceInterface) *models.User {
	// Get session cookie using the session name from auth service
	cookie, err := r.Cookie(authService.GetSessionName())
	if err != nil {
		// No session cookie found
		return nil
	}

	// Validate session with backend
	user, err := authService.ValidateSession(cookie.Value)
	if err != nil {
		// Session is invalid or expired
		return nil
	}

	return user
}

// GetSessionCookie is a helper function to get the session cookie by name from auth service
func GetSessionCookie(r *http.Request, authService AuthServiceInterface) (*http.Cookie, error) {
	return r.Cookie(authService.GetSessionName())
}

// ValidateUserSession - Get user from session, return AuthError if not found
func ValidateUserSession(r *http.Request, authService AuthServiceInterface) (*models.User, *http.Cookie, error) {
	user := GetUserFromSession(r, authService)
	if user == nil {
		return nil, nil, NewAuthError("Please log in to continue")
	}

	sessionCookie, err := GetSessionCookie(r, authService)
	if err != nil {
		return nil, nil, NewAuthError("Invalid session - please log in again")
	}

	return user, sessionCookie, nil
}

// GetOptionalUserSession - Get user from session, but don't return error if not found
func GetOptionalUserSession(r *http.Request, authService AuthServiceInterface) (*models.User, *http.Cookie) {
	user := GetUserFromSession(r, authService)
	if user == nil {
		return nil, nil
	}

	sessionCookie, err := GetSessionCookie(r, authService)
	if err != nil {
		return user, nil // Return user but no cookie if cookie retrieval fails
	}

	return user, sessionCookie
}

// AuthServiceInterface defines what we need from AuthService to avoid circular imports
type AuthServiceInterface interface {
	GetSessionName() string
	ValidateSession(sessionID string) (*models.User, error)
}
