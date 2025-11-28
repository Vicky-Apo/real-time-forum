package handlers

import (
	"frontend-service/internal/models"
	"frontend-service/internal/services"
	"frontend-service/internal/utils"
	"net/http"
)

type CategoryHandler struct {
	authService     *services.AuthService
	postService     *services.PostService
	categoryService *services.CategoryService
	templateService *services.TemplateService
	errorHandler    *SimpleErrorHandler
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(authService *services.AuthService, postService *services.PostService, categoryService *services.CategoryService, templateService *services.TemplateService) *CategoryHandler {
	return &CategoryHandler{
		authService:     authService,
		postService:     postService,
		categoryService: categoryService,
		templateService: templateService,
		errorHandler:    NewSimpleErrorHandler(templateService),
	}
}

func (h *CategoryHandler) ServeCategoryPosts(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user
	user := utils.GetUserFromSession(r, h.authService)

	// Extract category ID from URL
	categoryID := r.PathValue("id")
	if categoryID == "" {
		h.errorHandler.ShowError(w, "Category Not Found", "The requested category does not exist.")
		return
	}

	// 🔧 FIX: Use backend-style pagination utilities
	limit, offset := utils.ParsePaginationParams(r)
	sortBy := utils.ParseSortFromRequest(r, "newest")
	if sortBy != "newest" && sortBy != "oldest" && sortBy != "likes" && sortBy != "comments" {
		h.errorHandler.ShowError(w, "Invalid Sort Option", "The requested sort option is not valid. Please use 'newest', 'oldest', 'likes', or 'comments'.")
		return
	}

	// Get session cookie for API requests
	sessionCookie, _ := utils.GetSessionCookie(r, h.authService)

	// 🔧 FIXED: Now GetPostsByCategory returns PaginatedPostsResponse
	postsResponse, err := h.postService.GetPostsByCategory(categoryID, limit, offset, sortBy, sessionCookie)
	if err != nil {
		h.errorHandler.ShowError(w, "Failed to Load Posts", "We're having trouble loading the posts right now. Please try again later.")
		return
	}

	// Get all categories for sidebar
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		categories = []models.Category{} // Empty fallback
	}

	// 🔧 SIMPLIFIED: Get current category info (more reliable approach)
	var currentCategory *models.Category

	// First try to get category info from the first post (if posts exist)
	if len(postsResponse.Posts) > 0 && len(postsResponse.Posts[0].Categories) > 0 {
		firstPostCategory := postsResponse.Posts[0].Categories[0]
		currentCategory = &models.Category{
			ID:    firstPostCategory.ID,
			Name:  firstPostCategory.Name,
			Count: postsResponse.Pagination.TotalCount,
		}
	} else {
		// Fallback: Find category by ID from all categories list
		for _, cat := range categories {
			if cat.ID == categoryID {
				currentCategory = &models.Category{
					ID:    cat.ID,
					Name:  cat.Name,
					Count: 0, // No posts in this category
				}
				break
			}
		}
	}

	// Handle case where category not found
	if currentCategory == nil {
		h.errorHandler.ShowError(w, "Category Not Found", "The requested category does not exist.")
		return
	}

	// 🔧 FIXED: Use EXACT same data structure as working home and user-posts handlers
	data := map[string]interface{}{
		"Category":   currentCategory,          // Category info for page title and sidebar
		"Posts":      postsResponse.Posts,      // Direct array like other handlers
		"Pagination": postsResponse.Pagination, // Direct pagination like other handlers
		"Categories": categories,               // All categories for sidebar
		"User":       user,                     // Current user
		"Sort":       sortBy,                   // 🔧 ADDED: Sort parameter for pagination URLs
	}

	// Render template
	if err := h.templateService.Render(w, "category.html", data); err != nil {
		h.errorHandler.ShowError(w, "Failed to Render Page", "We're having trouble rendering the page right now. Please try again later.")
	}
}
