package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"real-time-forum/internal/middleware"
	"real-time-forum/internal/models"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/utils"
)
// ToggleCommentReactionHandler handles toggling reactions on comments
func ToggleCommentReactionHandler(crr *repository.CommentReactionRepository, nr *repository.NotificationRepository, cr *repository.CommentRepository, ur *repository.UserRepository, pr *repository.PostsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)
		if user == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// Parse request body
		var req models.CommentReactionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Validate the request
		if err := req.Validate(); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Call repository method to toggle comment reaction
		result, err := crr.ToggleCommentReaction(user.ID, req.CommentID, req.ReactionType)
		if err != nil {
			// Handle specific errors from repository
			if err.Error() == "comment not found" {
				utils.RespondWithError(w, http.StatusNotFound, "Comment not found")
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to toggle reaction")
			return
		}

		// Create notification only for new reactions
		if result.Action == models.ActionCommentLikeCreated || result.Action == models.ActionCommentDislikeCreated {
			// Get comment details to know who to notify
			comment, err := cr.GetCommentByID(req.CommentID, nil)
			if err == nil && comment.UserID != user.ID { // Don't notify yourself
				// Get post content preview
				post, err := pr.GetPostByID(comment.PostID, user.ID)
				if err == nil {
					contentPreview := post.Content
					if len(contentPreview) > 50 {
						contentPreview = contentPreview[:50] + "..."
					}

					// ✅ FIXED: Complete message for COMMENT reactions
					var actionText string
					if req.ReactionType == models.ReactionTypeLike {
						actionText = "liked your comment on post"
					} else {
						actionText = "disliked your comment on post"
					}

					// Create notification
					notification := &models.Notification{
						NotificationID:     utils.GenerateUUIDToken(),
						UserID:             comment.UserID, // Notify comment owner
						TriggerUsername:    user.Username,
						PostContentPreview: contentPreview,
						PostID:             comment.PostID,
						Action:             actionText, // ✅ Now says "liked your comment on"
						IsRead:             false,
						CreatedAt:          time.Now(),
					}

					// Save notification
					nr.CreateNotification(notification)
				}
			}
		}

		// Return the detailed reaction result
		utils.RespondWithSuccess(w, http.StatusOK, result)
	}
}
