// internal/services/comment_service.go
package services

import (
	"frontend-service/internal/models"
	"frontend-service/internal/utils"
	"net/http"
)

type CommentService struct {
	*BaseClient
}

// NewCommentService creates a new comment service
func NewCommentService(baseClient *BaseClient) *CommentService {
	return &CommentService{
		BaseClient: baseClient,
	}
}

// CreateComment creates a new comment on a post
func (s *CommentService) CreateComment(postID, content string, sessionCookie *http.Cookie) (*models.Comment, error) {
	// Prepare request data
	requestData := map[string]interface{}{
		"content": content,
	}

	// Make API request using utils
	apiResponse, err := utils.MakePOSTRequest(s.HTTPClient, s.BaseURL, "/comments/create-on-post/"+postID, requestData, sessionCookie)
	if err != nil {
		return nil, err
	}

	// Convert response to comment
	var comment models.Comment
	if err := utils.ConvertAPIData(apiResponse.Data, &comment); err != nil {
		return nil, utils.NewGeneralError("Failed to parse comment data", 500)
	}

	return &comment, nil
}

// UpdateComment updates an existing comment
func (s *CommentService) UpdateComment(commentID, content string, sessionCookie *http.Cookie) error {
	// Prepare request data
	requestData := map[string]interface{}{
		"content": content,
	}

	// Make API request using utils
	_, err := utils.MakePUTRequest(s.HTTPClient, s.BaseURL, "/comments/edit/"+commentID, requestData, sessionCookie)
	return err
}

// DeleteComment deletes a comment
func (s *CommentService) DeleteComment(commentID string, sessionCookie *http.Cookie) error {
	// Make API request using utils
	_, err := utils.MakeDELETERequest(s.HTTPClient, s.BaseURL, "/comments/remove/"+commentID, sessionCookie)
	return err
}

// GetCommentByID retrieves a single comment by ID
func (s *CommentService) GetCommentByID(commentID string, sessionCookie *http.Cookie) (*models.Comment, error) {
	// Make API request using utils
	apiResponse, err := utils.MakeGETRequest(s.HTTPClient, s.BaseURL, "/comments/view/"+commentID, sessionCookie)
	if err != nil {
		return nil, err
	}

	// Convert response to comment
	var comment models.Comment
	if err := utils.ConvertAPIData(apiResponse.Data, &comment); err != nil {
		return nil, utils.NewGeneralError("Failed to parse comment data", 500)
	}

	return &comment, nil
}
