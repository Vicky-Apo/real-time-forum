package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"platform.zone01.gr/git/gpapadopoulos/forum/internal/middleware"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/models"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/repository"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/utils"
)

func TogglePostReactionHandler(prr *repository.PostReactionRepository, nr *repository.NotificationRepository, pr *repository.PostsRepository, ur *repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// Parse request body
		var req models.PostReactionRequest
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

		// Call repository method to toggle post reaction
		result, err := prr.TogglePostReaction(user.ID, req.PostID, req.ReactionType)
		if err != nil {
			// Handle specific errors from repository
			if err.Error() == "post not found" {
				utils.RespondWithError(w, http.StatusNotFound, "Post not found")
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to toggle reaction")
			return
		}

		// Create notification only for new reactions
		if result.Action == models.ActionPostLikeCreated || result.Action == models.ActionPostDislikeCreated {
			// Get post details to know who to notify
			post, err := pr.GetPostByID(req.PostID, user.ID)
			if err == nil && post.UserID != user.ID { // Don't notify yourself
				// Get post content preview (first 50 chars)
				contentPreview := post.Content
				if len(contentPreview) > 50 {
					contentPreview = contentPreview[:50] + "..."
				}

				// ✅ FIXED: Use clean action text
				var actionText string
				if req.ReactionType == models.ReactionTypeLike {
					actionText = "liked your post"
				} else {
					actionText = "disliked your post"
				}

				// Create notification
				notification := &models.Notification{
					NotificationID:     utils.GenerateUUIDToken(),
					UserID:             post.UserID,
					TriggerUsername:    user.Username,
					PostContentPreview: contentPreview,
					PostID:             req.PostID,
					Action:             actionText, // ✅ Clean text instead of result.Action
					IsRead:             false,
					CreatedAt:          time.Now(),
				}

				// Save notification
				nr.CreateNotification(notification)
			}
		}

		// Return the detailed reaction result
		utils.RespondWithSuccess(w, http.StatusOK, result)
	}
}
