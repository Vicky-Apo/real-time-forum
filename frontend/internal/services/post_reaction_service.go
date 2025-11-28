// internal/services/post_reaction_service.go
package services

import (
	"net/http"

	"frontend-service/internal/models"
	"frontend-service/internal/utils"
)

type PostReactionService struct {
	*BaseClient
}

// NewPostReactionService creates a new post reaction service
func NewPostReactionService(baseClient *BaseClient) *PostReactionService {
	return &PostReactionService{
		BaseClient: baseClient,
	}
}

// TogglePostReaction toggles a like/dislike reaction on a post
func (s *PostReactionService) TogglePostReaction(postID string, reactionType int, sessionCookie *http.Cookie) (*models.ReactionResult, error) {
	// Prepare and validate request data
	requestData := models.PostReactionRequest{
		PostID:       postID,
		ReactionType: reactionType,
	}

	// Basic validation
	if err := s.validateReactionRequest(requestData); err != nil {
		return nil, err
	}

	// Make API request using utils
	apiResponse, err := utils.MakePOSTRequest(s.HTTPClient, s.BaseURL, "/reactions/posts/toggle", requestData, sessionCookie)
	if err != nil {
		return nil, err
	}

	// Convert response to reaction result
	var reactionResult models.ReactionResult
	if err := utils.ConvertAPIData(apiResponse.Data, &reactionResult); err != nil {
		return nil, utils.NewGeneralError("Failed to parse reaction result", 500)
	}

	return &reactionResult, nil
}

// GetPostReactionStatus gets the current user's reaction status on a post
func (s *PostReactionService) GetPostReactionStatus(postID string, sessionCookie *http.Cookie) (*int, error) {
	// Make API request using utils
	apiResponse, err := utils.MakeGETRequest(s.HTTPClient, s.BaseURL, "/reactions/posts/status/"+postID, sessionCookie)
	if err != nil {
		return nil, err
	}

	// Convert response to reaction status
	var statusResponse struct {
		PostID       string `json:"post_id"`
		UserReaction *int   `json:"user_reaction"` // nil, 1=like, 2=dislike
	}

	if err := utils.ConvertAPIData(apiResponse.Data, &statusResponse); err != nil {
		return nil, utils.NewGeneralError("Failed to parse reaction status", 500)
	}

	return statusResponse.UserReaction, nil
}

// validateReactionRequest validates post reaction request data
func (s *PostReactionService) validateReactionRequest(requestData models.PostReactionRequest) error {
	if requestData.PostID == "" {
		return utils.NewGeneralError("Post ID is required", 400)
	}

	if requestData.ReactionType != models.ReactionTypeLike && requestData.ReactionType != models.ReactionTypeDislike {
		return utils.NewGeneralError("Invalid reaction type: must be 1 (like) or 2 (dislike)", 400)
	}

	return nil
}
