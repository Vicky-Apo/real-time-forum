// simple_error.go
package handlers

import (
	"frontend-service/internal/services"
	"net/http"
)

// SimpleError represents a basic error
type SimpleError struct {
	Title   string
	Message string
}

// SimpleErrorHandler handles errors in the simplest way possible
type SimpleErrorHandler struct {
	templateService *services.TemplateService
}

// NewSimpleErrorHandler creates a simple error handler
func NewSimpleErrorHandler(templateService *services.TemplateService) *SimpleErrorHandler {
	return &SimpleErrorHandler{
		templateService: templateService,
	}
}

// ShowError displays an error page with just title, message, and home button
func (h *SimpleErrorHandler) ShowError(w http.ResponseWriter, title, message string) {
	data := SimpleError{
		Title:   title,
		Message: message,
	}

	if err := h.templateService.Render(w, "simple_error.html", data); err != nil {
		// If template fails, just log and send basic response
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred. Please go back to the <a href='/'>home page</a>."))
	}
}

func (h *SimpleErrorHandler) ParseAuthError(w http.ResponseWriter, title, message string) {
	data := SimpleError{
		Title:   title,
		Message: message,
	}

	if err := h.templateService.Render(w, "auth_error.html", data); err != nil {
		// If template fails, just log and send basic response
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unauthorized. Please login to <a href='/'>login page</a>."))
	}
}

func (h *SimpleErrorHandler) ShowAuthError(w http.ResponseWriter, r *http.Request) {
	h.ParseAuthError(w, "Authentication Required", "Please log in to continue.")
}

func (h *SimpleErrorHandler) ShowOAuthError(w http.ResponseWriter, message string) {
	h.ShowError(w, "Login Error", message)
}
