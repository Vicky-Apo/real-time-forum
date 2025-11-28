package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"real-time-forum/config"
	"real-time-forum/internal/middleware"
	"real-time-forum/internal/models"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/utils"
)

// CreatePostHandler handles creating a new post
func CreatePostHandler(pr repository.PostsRepositoryInterface, cr repository.CategoryRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		var req models.CreatePostRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if utils.ValidatePostContent(req.Content) != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid post content")
			return
		}

		// Validate and get category IDs
		var categoryIDs []string
		for _, categoryName := range req.CategoryNames {
			categoryID, err := cr.GetCategoryID(categoryName)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid category: "+categoryName)
				return
			}
			categoryIDs = append(categoryIDs, categoryID)
		}
		if len(categoryIDs) < config.Config.MinCategories {
			utils.RespondWithError(w, http.StatusBadRequest,
				fmt.Sprintf("Minimum %d category required", config.Config.MinCategories))
			return
		}

		if len(categoryIDs) > config.Config.MaxCategories {
			utils.RespondWithError(w, http.StatusBadRequest,
				fmt.Sprintf("Maximum %d categories allowed", config.Config.MaxCategories))
			return
		}
		// Create post - now returns lightweight response
		createResponse, err := pr.CreatePost(user.ID, req.Content, categoryIDs)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create post")
			return
		}

		// Return lightweight response
		utils.RespondWithSuccess(w, http.StatusCreated, createResponse)
	}
}

// GetAllPostsHandler retrieves all posts with pagination and sorting
func GetAllPostsHandler(pr repository.PostsRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// Parse pagination parameters - ONE LINE!
		limit, offset := utils.ParsePaginationParams(r)
		// Parse sort options from query parameters
		sortOptions := utils.ParsePostSortOptions(r)
		// Get posts and total count
		posts, err := pr.GetAllPosts(limit, offset, user.ID, sortOptions)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve posts")
			return
		}

		totalCount, err := pr.GetCountTotalPosts()
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve posts count")
			return
		}

		// Respond with standardized format - ONE LINE!
		utils.RespondWithPaginatedPosts(w, posts, totalCount, limit, offset)
	}
}

// GetPostsByCategoryHandler retrieves posts filtered by category with pagination and sorting
func GetPostsByCategoryHandler(pr repository.PostsRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// Extract category ID from URL path
		categoryID := r.PathValue("id")
		if categoryID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Category ID is required")
			return
		}
		// Parse pagination parameters
		limit, offset := utils.ParsePaginationParams(r)
		// Parse sort options from query parameters
		sortOptions := utils.ParsePostSortOptions(r)
		// Get total count for this category first
		totalCount, err := pr.GetCountPostByCategory(categoryID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve posts count")
			return
		}

		// Pass userID to repository
		posts, err := pr.GetPostsByCategory(categoryID, limit, offset, user.ID, sortOptions)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve posts")
			return
		}
		utils.RespondWithPaginatedPosts(w, posts, totalCount, limit, offset)
	}
}

// GetSinglePostHandler retrieves a single post with full details and reactions
func GetSinglePostHandler(pr repository.PostsRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// Extract post ID from URL path
		postID := r.PathValue("id")
		if postID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Post ID is required")
			return
		}

		//Pass userID to repository
		post, err := pr.GetPostByID(postID, user.ID)
		if err != nil {
			if err.Error() == "post not found" {
				utils.RespondWithError(w, http.StatusNotFound, "Post not found")
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve post")
			return
		}
		utils.RespondWithSuccess(w, http.StatusOK, post)
	}
}
