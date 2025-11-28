// internal/services/comment_reaction_service.go
package services

import (
	"net/http"

	"frontend-service/internal/models"
	"frontend-service/internal/utils"
)

type CommentReactionService struct {
	*BaseClient
}

// NewCommentReactionService creates a new comment reaction service
func NewCommentReactionService(baseClient *BaseClient) *CommentReactionService {
	return &CommentReactionService{
		BaseClient: baseClient,
	}
}

// ToggleCommentReaction toggles a like/dislike reaction on a comment
func (s *CommentReactionService) ToggleCommentReaction(commentID string, reactionType int, sessionCookie *http.Cookie) (*models.ReactionResult, error) {
	// Prepare and validate request data
	requestData := models.CommentReactionRequest{
		CommentID:    commentID,
		ReactionType: reactionType,
	}

	// Basic validation
	if err := s.validateReactionRequest(requestData); err != nil {
		return nil, err
	}

	// Make API request using utils
	apiResponse, err := utils.MakePOSTRequest(s.HTTPClient, s.BaseURL, "/reactions/comments/toggle", requestData, sessionCookie)
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

// GetCommentReactionStatus gets the current user's reaction status on a comment
func (s *CommentReactionService) GetCommentReactionStatus(commentID string, sessionCookie *http.Cookie) (*int, error) {
	// Make API request using utils
	apiResponse, err := utils.MakeGETRequest(s.HTTPClient, s.BaseURL, "/reactions/comments/status/"+commentID, sessionCookie)
	if err != nil {
		return nil, err
	}

	// Convert response to reaction status
	var reactionStatus struct {
		UserReaction *int `json:"user_reaction"`
	}

	if err := utils.ConvertAPIData(apiResponse.Data, &reactionStatus); err != nil {
		return nil, utils.NewGeneralError("Failed to parse reaction status", 500)
	}

	return reactionStatus.UserReaction, nil
}

// validateReactionRequest validates comment reaction request data
func (s *CommentReactionService) validateReactionRequest(requestData models.CommentReactionRequest) error {
	if requestData.CommentID == "" {
		return utils.NewGeneralError("Comment ID is required", 400)
	}

	if requestData.ReactionType != models.ReactionTypeLike && requestData.ReactionType != models.ReactionTypeDislike {
		return utils.NewGeneralError("Invalid reaction type: must be 1 (like) or 2 (dislike)", 400)
	}

	return nil
}
