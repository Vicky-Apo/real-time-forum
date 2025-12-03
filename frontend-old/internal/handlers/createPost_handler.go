package handlers

import (
	"net/http"

	"frontend-service/internal/services"
	"frontend-service/internal/utils"
)

type CreatePostHandler struct {
	authService     *services.AuthService
	postService     *services.PostService
	categoryService *services.CategoryService
	templateService *services.TemplateService
	errorHandler    *SimpleErrorHandler
}

// NewCreatePostHandler creates a new create post handler
func NewCreatePostHandler(authService *services.AuthService, postService *services.PostService, categoryService *services.CategoryService, templateService *services.TemplateService) *CreatePostHandler {
	return &CreatePostHandler{
		authService:     authService,
		postService:     postService,
		categoryService: categoryService,
		templateService: templateService,
		errorHandler:    NewSimpleErrorHandler(templateService),
	}
}

// ServeCreatePost handles both GET (show form) and POST (submit form) for creating posts
func (h *CreatePostHandler) ServeCreatePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.showCreatePostForm(w, r)
	case http.MethodPost:
		h.handleCreatePostForm(w, r)
	default:
		h.errorHandler.ShowError(w, "Method Not Allowed", "This method is not allowed for creating posts.")
	}
}

// showCreatePostForm displays the create post form (GET request)
func (h *CreatePostHandler) showCreatePostForm(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	user := utils.GetUserFromSession(r, h.authService)
	if user == nil {
		h.errorHandler.ShowAuthError(w, r)
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
		"Categories": categories,
		"FormData":   map[string]interface{}{}, // Empty form data for initial load
	}

	// Render the template
	if err := h.templateService.Render(w, "create-post.html", data); err != nil {
		h.errorHandler.ShowError(w, "Failed to Render Page", "We're having trouble rendering the page right now. Please try again later.")
		return
	}
}

// handleCreatePostForm processes the create post form submission (POST request)
func (h *CreatePostHandler) handleCreatePostForm(w http.ResponseWriter, r *http.Request) {
	// 1. Check if user is logged in
	user := utils.GetUserFromSession(r, h.authService)
	if user == nil {
		h.errorHandler.ShowAuthError(w, r)
		return
	}

	// 2. Get the session cookie for backend authentication
	sessionCookie, err := utils.GetSessionCookie(r, h.authService)
	if err != nil {
		h.errorHandler.ShowAuthError(w, r)
		return
	}

	// 3. Call the service to forward the request and get the result
	postID, apiError, status := h.postService.ForwardCreatePost(r, sessionCookie)

	// 4. Handle any API/backend errors
	if apiError != "" {
		if status == http.StatusUnauthorized {
			h.errorHandler.ShowAuthError(w, r)
			return
		}
		h.showCreatePostError(w, r, apiError, nil)
		return
	}

	// 5. Success: redirect to the newly created post
	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}

// showCreatePostError displays the create post form with error message and preserved form data
func (h *CreatePostHandler) showCreatePostError(w http.ResponseWriter, r *http.Request, errorMsg string, formData map[string]interface{}) {
	// Get user (should be logged in if we reach this point)
	user := utils.GetUserFromSession(r, h.authService)

	// Get all categories for the form
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		h.errorHandler.ShowError(w, "Failed to Load Categories", "We're having trouble loading the categories right now. Please try again later.")
		return
	}

	// Prepare data for template
	data := map[string]interface{}{
		"User":       user,
		"Categories": categories,
		"Error":      errorMsg,
		"FormData":   formData,
	}

	// Render the template with error
	if err := h.templateService.Render(w, "create-post.html", data); err != nil {
		h.errorHandler.ShowError(w, "Failed to Render Page", "We're having trouble rendering the page right now. Please try again later.")
		return
	}
}
