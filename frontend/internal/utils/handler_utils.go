// internal/utils/handler_utils.go
package utils

import (
	"net/http"
	"net/url"
)

// HandleError - Route error to appropriate handler based on error type
func HandleError(w http.ResponseWriter, r *http.Request, err error, errorHandler ErrorHandler) {
	if appErr, ok := err.(*AppError); ok {
		switch appErr.Type {
		case AuthError:
			errorHandler.ShowAuthError(w, r)
		case HomeError:
			RedirectToHomeWithError(w, r, appErr.Message)
		case GeneralError:
			errorHandler.ShowError(w, appErr.Message, getStatusText(appErr.Code))
		}
	} else {
		// Fallback for unknown errors
		errorHandler.ShowError(w, "An unexpected error occurred", "Please try again later")
	}
}

// ErrorHandler interface for error handling
type ErrorHandler interface {
	ShowError(w http.ResponseWriter, title, message string)
	ShowAuthError(w http.ResponseWriter, r *http.Request)
}

// RedirectToHomeWithError - Redirect to home page with error message
func RedirectToHomeWithError(w http.ResponseWriter, r *http.Request, message string) {
	http.Redirect(w, r, "/?error="+url.QueryEscape(message), http.StatusSeeOther)
}

// getStatusText - Convert HTTP status code to readable text
func getStatusText(code int) string {
	switch code {
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusInternalServerError:
		return "Internal Server Error"
	case http.StatusBadGateway:
		return "Service Unavailable"
	default:
		return "Error"
	}
}
