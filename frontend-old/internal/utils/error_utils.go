// internal/utils/error_utils.go
package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"frontend-service/internal/models"
)

// Custom error types for different redirect behaviors
type AppError struct {
	Type    ErrorType
	Message string
	Code    int
}

type ErrorType int

const (
	// HomeError - Redirect to home page (post not found, invalid resource, etc.)
	HomeError ErrorType = iota
	// AuthError - Redirect to login page (need authentication)
	AuthError
	// GeneralError - Show error page (validation errors, server errors)
	GeneralError
)

func (e *AppError) Error() string {
	return e.Message
}

// NewHomeError - Creates error that redirects to home page
func NewHomeError(message string) *AppError {
	return &AppError{Type: HomeError, Message: message, Code: http.StatusNotFound}
}

// NewAuthError - Creates error that redirects to login page
func NewAuthError(message string) *AppError {
	return &AppError{Type: AuthError, Message: message, Code: http.StatusUnauthorized}
}

// NewGeneralError - Creates error that shows error page
func NewGeneralError(message string, code int) *AppError {
	return &AppError{Type: GeneralError, Message: message, Code: code}
}

// HandleHTTPStatus - Convert HTTP status codes to appropriate error types
func HandleHTTPStatus(statusCode int, body []byte) error {
	switch statusCode {
	case http.StatusOK, http.StatusCreated:
		return nil

	// AUTH ERRORS - Redirect to login
	case http.StatusUnauthorized:
		return NewAuthError("Please log in to continue")
	case http.StatusForbidden:
		if containsErrorMessage(body, "please log in") {
			return NewAuthError("Please log in to continue")
		}
		return NewGeneralError("Access denied - insufficient permissions", http.StatusForbidden)

	// HOME ERRORS - Redirect to home page
	case http.StatusNotFound:
		return NewHomeError("The requested resource was not found")

	// GENERAL ERRORS - Show error page
	case http.StatusBadRequest:
		message := extractErrorMessage(body, "Invalid request")
		return NewGeneralError(message, http.StatusBadRequest)
	case http.StatusInternalServerError:
		return NewGeneralError("Internal server error - please try again later", http.StatusInternalServerError)
	case http.StatusBadGateway:
		return NewGeneralError("Service temporarily unavailable", http.StatusBadGateway)
	default:
		return NewGeneralError(fmt.Sprintf("Unexpected error (code: %d)", statusCode), statusCode)
	}
}

// extractErrorMessage - Extract error message from API response
func extractErrorMessage(body []byte, defaultMessage string) string {
	var apiResponse models.APIResponse
	if json.Unmarshal(body, &apiResponse) == nil && apiResponse.Error != "" {
		return apiResponse.Error
	}
	return defaultMessage
}

// containsErrorMessage - Check if response contains specific error message
func containsErrorMessage(body []byte, text string) bool {
	var apiResponse models.APIResponse
	if json.Unmarshal(body, &apiResponse) == nil {
		return stringContains(apiResponse.Error, text)
	}
	return stringContains(string(body), text)
}

// stringContains - Simple string contains check
func stringContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr)))
}
