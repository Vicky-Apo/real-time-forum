package handlers

import (
	"frontend-service/internal/models"
	"frontend-service/internal/services"
	"frontend-service/internal/utils"
	"frontend-service/internal/validations"
	"net/http"
	"strings"
)

type CommentHandler struct {
	authService     *services.AuthService
	commentService  *services.CommentService
	postService     *services.PostService
	templateService *services.TemplateService
	errorHandler    *SimpleErrorHandler
}

// NewCommentHandler creates a new comment handler
func NewCommentHandler(authService *services.AuthService, commentService *services.CommentService, postService *services.PostService, templateService *services.TemplateService) *CommentHandler {
	return &CommentHandler{
		authService:     authService,
		commentService:  commentService,
		postService:     postService,
		templateService: templateService,
		errorHandler:    NewSimpleErrorHandler(templateService),
	}
}

// ServeCreateComment handles comment creation (form submission)
func (h *CommentHandler) ServeCreateComment(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		h.errorHandler.ShowError(w, "Method Not Allowed", "Only POST method is allowed.")
		return
	}

	// Check if user is logged in
	user := utils.GetUserFromSession(r, h.authService)
	if user == nil {
		h.errorHandler.ShowAuthError(w, r)
		return
	}

	// Extract post ID from URL path
	postID := r.PathValue("post_id")
	if postID == "" {
		h.errorHandler.ShowError(w, "Post ID Required", "Post ID is required to create a comment.")
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/post/"+postID+"?error=invalid_form", http.StatusSeeOther)
		return
	}

	// Get and validate comment content
	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Redirect(w, r, "/post/"+postID+"?error=empty_content", http.StatusSeeOther)
		return
	}

	if err := validations.ValidateCommentContent(content); err != nil {
		http.Redirect(w, r, "/post/"+postID+"?error=validation_failed", http.StatusSeeOther)
		return
	}

	// Get session cookie for API call
	sessionCookie, err := utils.GetSessionCookie(r, h.authService) // CHANGED: Use utility function instead of hardcoded "session_id"
	if err != nil {
		h.errorHandler.ShowAuthError(w, r)
		return
	}

	// Create comment via API
	_, err = h.commentService.CreateComment(postID, content, sessionCookie)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			h.errorHandler.ShowAuthError(w, r)
			return
		}
		http.Redirect(w, r, "/post/"+postID+"?error=create_failed", http.StatusSeeOther)
		return
	}

	// Redirect back to the post (comment will appear after page refresh)
	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}

// ServeEditComment handles both GET (show edit form) and POST (save changes)
func (h *CommentHandler) ServeEditComment(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	user := utils.GetUserFromSession(r, h.authService)
	if user == nil {
		h.errorHandler.ShowAuthError(w, r)
		return
	}

	// Extract comment ID from URL path
	commentID := r.PathValue("comment_id")
	if commentID == "" {
		h.errorHandler.ShowError(w, "Comment ID Required", "Comment ID is required to edit a comment.")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.showEditCommentForm(w)

	case http.MethodPost:
		h.handleEditCommentForm(w, r, commentID)
	default:
		h.errorHandler.ShowError(w, "Method Not Allowed", "Only POST method is allowed.")
	}
}

// ServeDeleteComment handles comment deletion (form submission)
func (h *CommentHandler) ServeDeleteComment(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		h.errorHandler.ShowError(w, "Method Not Allowed", "Only POST method is allowed.")
		return
	}

	// Check if user is logged in
	user := utils.GetUserFromSession(r, h.authService)
	if user == nil {
		h.errorHandler.ShowAuthError(w, r)
		return
	}

	// Extract comment ID from URL path
	commentID := r.PathValue("comment_id")
	if commentID == "" {
		h.errorHandler.ShowError(w, "Comment ID Required", "Comment ID is required to edit a comment.")
		return
	}

	// Get session cookie for API call
	sessionCookie, err := utils.GetSessionCookie(r, h.authService) // CHANGED: Use utility function instead of hardcoded "session_id"
	if err != nil {
		h.errorHandler.ShowAuthError(w, r)
		return
	}

	// Delete comment via API
	err = h.commentService.DeleteComment(commentID, sessionCookie)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			h.errorHandler.ShowAuthError(w, r)
			return
		}
		// We don't know the post ID here, so redirect to home with error
		http.Redirect(w, r, "/?error=delete_failed", http.StatusSeeOther)
		return
	}

	// Try to get the referring post URL to redirect back
	referer := r.Header.Get("Referer")
	if referer != "" && strings.Contains(referer, "/post/") {
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}

	// Fallback: redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// showEditCommentForm displays the edit comment form (GET request)
func (h *CommentHandler) showEditCommentForm(w http.ResponseWriter) {

	h.errorHandler.ShowError(w, "Not Implemented", "Edit comment form is not implemented yet. Please try again later.")
}

// handleEditCommentForm processes the edit comment form submission (POST request)
func (h *CommentHandler) handleEditCommentForm(w http.ResponseWriter, r *http.Request, commentID string) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		h.errorHandler.ShowError(w, "Invalid Form Data", "Please provide valid form data.")
		return
	}

	// Get and validate comment content
	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		h.errorHandler.ShowError(w, "Comment Content Required", "Comment content is required.")
		return
	}

	if err := validations.ValidateCommentContent(content); err != nil {
		h.errorHandler.ShowError(w, "Invalid Comment Content", err.Error())
		return
	}

	// Get session cookie for API call
	sessionCookie, err := utils.GetSessionCookie(r, h.authService) // Use utility function instead of hardcoded "session_id"
	if err != nil {
		h.errorHandler.ShowAuthError(w, r)
		return
	}

	// Update comment via API
	err = h.commentService.UpdateComment(commentID, content, sessionCookie)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			h.errorHandler.ShowAuthError(w, r)
			return
		}
		if strings.Contains(err.Error(), "forbidden") {
			h.errorHandler.ShowError(w, "Forbidden", "You can only edit your own comments.")
			return
		}
		if strings.Contains(err.Error(), "not found") {
			h.errorHandler.ShowError(w, "Comment Not Found", "The requested comment does not exist.")
			return
		}
		h.errorHandler.ShowError(w, "Failed to Update Comment", "We're having trouble updating the comment right now. Please try again later.")
		return
	}

	// Try to redirect back to the referring post
	referer := r.Header.Get("Referer")
	if referer != "" && strings.Contains(referer, "/post/") {
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}

	// Fallback: redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ServeEditCommentForm shows the edit comment form (GET)
func (h *CommentHandler) ServeEditCommentForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler.ShowError(w, "Method Not Allowed", "Only GET method is allowed.")
		return
	}

	// Get authenticated user
	user := utils.GetUserFromSession(r, h.authService)
	if user == nil {
		h.errorHandler.ShowAuthError(w, r)
		return
	}

	// Get comment ID from URL
	commentID := r.PathValue("id")
	if commentID == "" {
		h.errorHandler.ShowError(w, "Comment ID Required", "Comment ID is required to edit a comment.")
		return
	}

	// Get session cookie
	sessionCookie, _ := utils.GetSessionCookie(r, h.authService)

	// Get the comment
	comment, err := h.commentService.GetCommentByID(commentID, sessionCookie)
	if err != nil {
		h.errorHandler.ShowError(w, "Comment Not Found", "The requested comment does not exist.")
		return
	}

	// Check if user owns the comment
	if comment.UserID != user.ID {
		h.errorHandler.ShowError(w, "Forbidden", "You can only edit your own comments.")
		return
	}

	// Prepare template data
	data := struct {
		Comment *models.Comment
		User    *models.User
		Error   string
	}{
		Comment: comment,
		User:    user,
	}

	// Render edit comment template
	if err := h.templateService.Render(w, "edit-comment.html", data); err != nil {

		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return

	}
}

// ServeEditCommentSubmit processes the edit comment form (POST)
func (h *CommentHandler) ServeEditCommentSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorHandler.ShowError(w, "Method Not Allowed", "Only POST method is allowed.")
		return
	}

	// Get authenticated user
	user := utils.GetUserFromSession(r, h.authService)
	if user == nil {
		h.errorHandler.ShowAuthError(w, r)
		return
	}

	// Get comment ID from URL
	commentID := r.PathValue("id")
	if commentID == "" {
		h.errorHandler.ShowError(w, "Comment ID Required", "Comment ID is required to edit a comment.")
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		h.errorHandler.ShowError(w, "Invalid Form Data", "Please provide valid form data.")
		return
	}

	// Get form data
	content := strings.TrimSpace(r.FormValue("content"))
	redirectTo := r.FormValue("redirect_to")

	// Validate content
	if err := validations.ValidateCommentContent(content); err != nil {
		// Show form again with error
		h.showEditCommentError(w, r, commentID, content, err.Error())
		return
	}

	// Get session cookie
	sessionCookie, _ := utils.GetSessionCookie(r, h.authService)

	// Update the comment
	if err := h.commentService.UpdateComment(commentID, content, sessionCookie); err != nil {
		h.errorHandler.ShowError(w, "Failed to Update Comment", "We're having trouble updating the comment right now. Please try again later.")
		return
	}

	// Redirect back to post or default location
	if redirectTo == "" {
		redirectTo = "/"
	}
	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
}

// Helper method to show edit form with error
func (h *CommentHandler) showEditCommentError(w http.ResponseWriter, r *http.Request, commentID, content, errorMsg string) {
	user := utils.GetUserFromSession(r, h.authService)

	// Create a comment object with the form data
	comment := &models.Comment{
		ID:      commentID,
		Content: content,
	}

	data := struct {
		Comment *models.Comment
		User    *models.User
		Error   string
	}{
		Comment: comment,
		User:    user,
		Error:   errorMsg,
	}

	if err := h.templateService.Render(w, "edit-comment.html", data); err != nil {
		h.errorHandler.ShowError(w, "Failed to Render Page", "We're having trouble rendering the page right now. Please try again later.")
	}
}
