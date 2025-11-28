package handlers

import (
	"net/http"
	"strings"

	"frontend-service/internal/services"
	"frontend-service/internal/utils"
)

type DeletePostHandler struct {
	authService     *services.AuthService
	postService     *services.PostService
	errorHandler    *SimpleErrorHandler
	templateService *services.TemplateService
}

// NewDeletePostHandler creates a new delete post handler
func NewDeletePostHandler(authService *services.AuthService, postService *services.PostService, templateService *services.TemplateService) *DeletePostHandler {
	return &DeletePostHandler{
		authService:     authService,
		postService:     postService,
		errorHandler:    NewSimpleErrorHandler(templateService),
		templateService: templateService,
	}
}

// ServeDeletePost handles post deletion (POST request only for security)
func (h *DeletePostHandler) ServeDeletePost(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method for security (prevent accidental deletion via GET)
	if r.Method != http.MethodPost {
		h.errorHandler.ShowError(w, "Method Not Allowed", "This method is not allowed for deleting posts.")
		return
	}

	// Check if user is logged in
	user := utils.GetUserFromSession(r, h.authService)
	if user == nil {
		h.errorHandler.ShowAuthError(w, r)
		return
	}

	// Extract post ID from URL path
	postID := r.PathValue("id")
	if postID == "" {
		h.errorHandler.ShowError(w, "Post ID is Required", "Post ID is required to delete a post.")
		return
	}

	// Get the post to verify ownership
	sessionCookie, _ := utils.GetSessionCookie(r, h.authService) // CHANGED: Use utility function instead of hardcoded "session_id"
	post, _, err := h.postService.GetSinglePostWithComments(postID, 1, 0, "oldest", sessionCookie)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.errorHandler.ShowError(w, "Post Not Found", "The requested post does not exist.")
			return
		}
		h.errorHandler.ShowError(w, "Failed to Verify Post Ownership", "We're having trouble verifying post ownership right now. Please try again later.")
		return
	}

	// Check if user owns the post
	if post.UserID != user.ID {
		h.errorHandler.ShowError(w, "Forbidden", "You can only delete your own posts.")
		return
	}

	// Call backend API to delete post
	err = h.postService.DeletePost(postID, sessionCookie)
	if err != nil {
		// Handle different error types
		if strings.Contains(err.Error(), "unauthorized") {
			h.errorHandler.ShowAuthError(w, r)
			return
		}
		if strings.Contains(err.Error(), "forbidden") {
			h.errorHandler.ShowError(w, "Forbidden", "You can only delete your own posts.")
			return
		}
		if strings.Contains(err.Error(), "not found") {
			h.errorHandler.ShowError(w, "Post Not Found", "The requested post does not exist.")
			return
		}

		// For other errors, redirect back to post with error (we could implement flash messages later)
		http.Redirect(w, r, "/post/"+postID+"?error=delete_failed", http.StatusSeeOther)
		return
	}

	// Determine where to redirect after successful deletion
	redirectURL := "/"

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
