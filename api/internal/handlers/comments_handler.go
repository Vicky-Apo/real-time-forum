package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/internal/middleware"
	"real-time-forum/internal/models"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/utils"
)

// Create comment handler.
func CreateCommentHandler(cor *repository.CommentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Get authenticated user
		user := middleware.GetCurrentUser(r)
		if user == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

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

// Get ALL comments by post ID handler.
func GetCommentsByPostIDHandler(cor *repository.CommentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Get user context
		currentUser := middleware.GetCurrentUser(r)
		var userID *string = nil
		if currentUser != nil {
			userID = &currentUser.ID
		}

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
		comments, err := cor.GetCommentsByPostID(postID, limit, offset, userID, sortOptions)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve comments")
			return
		}

		// Respond with paginated comments
		utils.RespondWithPaginatedComments(w, comments, totalCount, limit, offset)
	}
}
