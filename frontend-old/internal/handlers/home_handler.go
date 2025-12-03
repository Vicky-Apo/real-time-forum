// home_handler.go
package handlers

import (
	"net/http"

	"frontend-service/internal/models"
	"frontend-service/internal/services"
	"frontend-service/internal/utils"
)

type HomeHandler struct {
	authService     *services.AuthService
	postService     *services.PostService
	categoryService *services.CategoryService
	templateService *services.TemplateService
	errorHandler    *SimpleErrorHandler
}

// NewHomeHandler creates a new home handler
func NewHomeHandler(authService *services.AuthService, postService *services.PostService, categoryService *services.CategoryService, templateService *services.TemplateService) *HomeHandler {
	return &HomeHandler{
		authService:     authService,
		postService:     postService,
		categoryService: categoryService,
		templateService: templateService,
		errorHandler:    NewSimpleErrorHandler(templateService),
	}
}

// ServeHome handles the main page request
func (h *HomeHandler) ServeHome(w http.ResponseWriter, r *http.Request) {
	// Parse pagination and sort parameters
	limit, offset := utils.ParsePaginationParams(r)
	sortBy := utils.ParseSortFromRequest(r, "newest")

	// Validate sort parameter
	if sortBy != "newest" && sortBy != "oldest" && sortBy != "likes" && sortBy != "comments" {
		h.errorHandler.ShowError(w, "Invalid Sort Option", "The selected sort option is not valid. Please use 'newest', 'oldest', 'likes', or 'comments'.")
		return
	}

	// Check if user is logged in and get session cookie
	user := utils.GetUserFromSession(r, h.authService)
	var sessionCookie *http.Cookie
	if user != nil {
		sessionCookie, _ = utils.GetSessionCookie(r, h.authService)
	}

	// Get posts from backend API
	postsResponse, err := h.postService.GetAllPosts(limit, offset, sortBy, sessionCookie)
	if err != nil {
		h.errorHandler.ShowError(w, "Unable to Load Posts", "We're having trouble loading the posts right now. Please try again later.")
		return
	}

	// Get categories from backend API
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		h.errorHandler.ShowError(w, "Unable to Load Categories", "We're having trouble loading the categories right now. Please try again later.")
		return
	}

	// Use structured HomePageData instead of generic map
	data := models.HomePageData{
		Posts:      postsResponse.Posts,
		Categories: categories,
		Pagination: postsResponse.Pagination,
		User:       user,
		Sort:       sortBy,
	}

	// Render the template
	if err := h.templateService.Render(w, "home.html", data); err != nil {
		h.errorHandler.ShowError(w, "Internal Server Error", "We're having trouble loading the home page right now. Please try again later.")
		return
	}
}
