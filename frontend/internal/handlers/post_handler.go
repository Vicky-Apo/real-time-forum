package handlers

import (
	"net/http"

	"frontend-service/internal/models"
	"frontend-service/internal/services"
	"frontend-service/internal/utils"
)

type PostHandler struct {
	authService     *services.AuthService
	postService     *services.PostService
	templateService *services.TemplateService
	errorHandler    *SimpleErrorHandler
}

// NewPostHandler creates a new post handler
func NewPostHandler(authService *services.AuthService, postService *services.PostService, templateService *services.TemplateService) *PostHandler {
	return &PostHandler{
		authService:     authService,
		postService:     postService,
		templateService: templateService,
		errorHandler:    NewSimpleErrorHandler(templateService),
	}
}

// ServePostView handles the single post view page request
func (h *PostHandler) ServePostView(w http.ResponseWriter, r *http.Request) {
	// Extract post ID from URL path
	postID := r.PathValue("id")
	if postID == "" {
		h.errorHandler.ShowError(w, "Post ID Required", "Please provide a valid post ID.")
		return
	}

	// 🔧 FIX: Use backend-style pagination utilities for comments
	limit, offset := utils.ParsePaginationParams(r)

	// Parse sort parameter for comments (default to newest for better UX)
	sortBy := utils.ParseSortFromRequest(r, "newest")

	// Validate sort parameter for comments
	if sortBy != "newest" && sortBy != "oldest" && sortBy != "likes" && sortBy != "comments" {
		h.errorHandler.ShowError(w, "Invalid Sort Option", "Please use 'newest', 'oldest', 'likes', or 'comments'.")
		return
	}

	// Check if user is logged in and get session cookie for reaction data
	user := utils.GetUserFromSession(r, h.authService)
	var sessionCookie *http.Cookie
	if user != nil {
		sessionCookie, _ = utils.GetSessionCookie(r, h.authService)
	}

	// 🔧 UPDATED: Get post and comments with pagination info
	post, commentsResponse, err := h.postService.GetSinglePostWithComments(postID, limit, offset, sortBy, sessionCookie)
	if err != nil {
		if err.Error() == "post not found" {
			h.errorHandler.ShowError(w, "Post Not Found", "The post you're looking for doesn't exist.")
			return
		}
		h.errorHandler.ShowError(w, "Failed to Load Post", "We're having trouble loading the post right now. Please try again later.")
		return
	}

	// 🔧 UPDATED: Convert []*Comment to []Comment for template compatibility
	var commentsSlice []models.Comment
	for _, comment := range commentsResponse.Comments {
		if comment != nil {
			commentsSlice = append(commentsSlice, *comment)
		}
	}

	// 🔧 UPDATED: Prepare data with comment pagination info
	data := models.PostPageData{
		Post:               post,
		Comments:           commentsSlice,
		CommentsPagination: commentsResponse.Pagination, // 🔧 ADDED: Comment pagination
		User:               user,
		Sort:               sortBy, // 🔧 ADDED: Current sort for template
	}

	// Render the template
	if err := h.templateService.Render(w, "post.html", data); err != nil {
		h.errorHandler.ShowError(w, "Failed to Render Post Page", "We're having trouble loading the post page right now. Please try again later.")
		return
	}
}
