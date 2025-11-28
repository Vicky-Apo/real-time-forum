package handlers

import (
	"net/http"
	"strings"

	"frontend-service/internal/services"
	"frontend-service/internal/utils"
)

type EditPostHandler struct {
	authService     *services.AuthService
	postService     *services.PostService
	categoryService *services.CategoryService
	templateService *services.TemplateService
	errorHandler    *SimpleErrorHandler
}

// NewEditPostHandler creates a new edit post handler
func NewEditPostHandler(authService *services.AuthService, postService *services.PostService, categoryService *services.CategoryService, templateService *services.TemplateService) *EditPostHandler {
	return &EditPostHandler{
		authService:     authService,
		postService:     postService,
		categoryService: categoryService,
		templateService: templateService,
		errorHandler:    NewSimpleErrorHandler(templateService),
	}
}

// ServeEditPost handles both GET (show form) and POST (submit form) for editing posts
func (h *EditPostHandler) ServeEditPost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.showEditPostForm(w, r)
	case http.MethodPost:
		h.handleEditPostForm(w, r)
	default:
		h.errorHandler.ShowError(w, "Method Not Allowed", "This method is not allowed for editing posts.")
	}
}

// showEditPostForm displays the edit post form (GET request)
func (h *EditPostHandler) showEditPostForm(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	user := utils.GetUserFromSession(r, h.authService)
	if user == nil {
		h.errorHandler.ShowAuthError(w, r)
		return
	}

	// Extract post ID from URL path
	postID := r.PathValue("id")
	if postID == "" {

		return
	}

	// Get the post to edit
	sessionCookie, _ := utils.GetSessionCookie(r, h.authService) // CHANGED: Use utility function instead of hardcoded "session_id"
	post, _, err := h.postService.GetSinglePostWithComments(postID, 1, 0, "oldest", sessionCookie)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.errorHandler.ShowError(w, "Post Not Found", "The requested post does not exist.")
			return
		}
		h.errorHandler.ShowError(w, "Failed to Load Post", "We're having trouble loading the post right now. Please try again later.")
		return
	}

	// Check if user owns the post
	if post.UserID != user.ID {
		h.errorHandler.ShowError(w, "Forbidden", "You can only edit your own posts.")
		return
	}

	// Get all categories for the form
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		h.errorHandler.ShowError(w, "Failed to Load Categories", "We're having trouble loading the categories right now. Please try again later.")
		return
	}

	// Extract category names from the post
	var postCategoryNames []string
	for _, cat := range post.Categories {
		postCategoryNames = append(postCategoryNames, cat.Name)
	}

	// Prepare data for template
	data := map[string]interface{}{
		"User":       user,
		"Post":       post,
		"Categories": categories,
		"FormData": map[string]interface{}{
			"content":    post.Content,
			"categories": postCategoryNames,
		},
		"IsEdit": true, // Flag to indicate this is an edit form
	}

	// Render the template (we'll reuse create-post.html with edit mode)
	if err := h.templateService.Render(w, "edit-post.html", data); err != nil {
		h.errorHandler.ShowError(w, "Failed to Render Page", "We're having trouble rendering the page right now. Please try again later.")
		return
	}
}

// handleEditPostForm processes the edit post form submission (POST request)
func (h *EditPostHandler) handleEditPostForm(w http.ResponseWriter, r *http.Request) {
	// 1. Check login
	user := utils.GetUserFromSession(r, h.authService)
	if user == nil {
		h.errorHandler.ShowAuthError(w, r)
		return
	}

	// 2. Extract post ID
	postID := r.PathValue("id")
	if postID == "" {
		h.errorHandler.ShowError(w, "Post Not Found", "The requested post does not exist.")
		return
	}

	// 3. Get session cookie for backend
	sessionCookie, err := utils.GetSessionCookie(r, h.authService)
	if err != nil {
		h.errorHandler.ShowAuthError(w, r)
		return
	}

	// 4. Forward the entire multipart/form-data request to backend for updating the post (including images/removals)
	err = h.postService.UpdatePost(r, postID, sessionCookie)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			h.errorHandler.ShowAuthError(w, r)
			return
		}
		if strings.Contains(err.Error(), "forbidden") {
			h.errorHandler.ShowError(w, "Forbidden", "You can only edit your own posts.")
			return
		}
		// For all other errors, show the error and re-render the edit form
		h.showEditPostError(w, r, postID, err.Error(), nil)
		return
	}

	// 5. Success: redirect to the updated post
	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}

// showEditPostError displays the edit post form with error message and preserved form data
func (h *EditPostHandler) showEditPostError(w http.ResponseWriter, r *http.Request, postID, errorMsg string, formData map[string]interface{}) {
	// Get user (should be logged in if we reach this point)
	user := utils.GetUserFromSession(r, h.authService)

	// Get the original post for reference
	sessionCookie, _ := utils.GetSessionCookie(r, h.authService) // CHANGED: Use utility function instead of hardcoded "session_id"
	post, _, err := h.postService.GetSinglePostWithComments(postID, 1, 0, "oldest", sessionCookie)
	if err != nil {
		h.errorHandler.ShowError(w, "Failed to Load Post", "We're having trouble loading the post right now. Please try again later.")
		return
	}

	// Get all categories for the form
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		h.errorHandler.ShowError(w, "Failed to Load Categories", "We're having trouble loading the categories right now. Please try again later.")
		return
	}

	// Prepare data for template
	data := map[string]interface{}{
		"User":       user,
		"Post":       post,
		"Categories": categories,
		"Error":      errorMsg,
		"FormData":   formData,
		"IsEdit":     true,
	}

	// Render the template with error
	if err := h.templateService.Render(w, "edit-post.html", data); err != nil {
		h.errorHandler.ShowError(w, "Failed to Render Page", "We're having trouble rendering the page right now. Please try again later.")
		return
	}
}
