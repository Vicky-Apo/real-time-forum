package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/internal/middleware"
	"real-time-forum/internal/models"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/utils"
)

// CreateCommentHandler handles creating a new comment on a post
func CreateCommentHandler(cor repository.CommentRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// Get postID from URL path, not request body
		postID := r.PathValue("id")
		if postID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Post ID is required")
			return
		}

		// Parse request body (only content, no postID)
		var req models.CreateCommentRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Validate comment content
		if err := utils.ValidateCommentContent(req.Content); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Create comment - now returns lightweight response
		createResponse, err := cor.CreateComment(postID, user.ID, req.Content)
		if err != nil {
			if err.Error() == "post not found" {
				utils.RespondWithError(w, http.StatusNotFound, "Post not found")
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create comment")
			return
		}

		// Return lightweight response
		utils.RespondWithSuccess(w, http.StatusCreated, createResponse)
	}
}

// GetCommentsByPostIDHandler retrieves all comments for a specific post
func GetCommentsByPostIDHandler(cor repository.CommentRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// Get post ID from URL
		postID := r.PathValue("id")
		if postID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Post ID is required")
			return
		}

		// Parse pagination and sort parameters using unified system
		limit, offset := utils.ParsePaginationParams(r)
		// Parse sort options using unified system
		sortOptions := utils.ParseCommentSortOptions(r)

		// Get total count
		totalCount, err := cor.GetCommentCountByPost(postID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve comment count")
			return
		}

		// Get comments with sorting
		comments, err := cor.GetCommentsByPostID(postID, limit, offset, user.ID, sortOptions)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve comments")
			return
		}

		// Respond with paginated comments
		utils.RespondWithPaginatedComments(w, comments, totalCount, limit, offset)
	}
}
