// internal/services/user_service.go
package services

import (
	"net/http"

	"frontend-service/internal/models"
	"frontend-service/internal/utils"
)

type UserService struct {
	*BaseClient
}

// NewUserService creates a new user service
func NewUserService(baseClient *BaseClient) *UserService {
	return &UserService{
		BaseClient: baseClient,
	}
}

// GetUserProfile retrieves user profile with stats from the backend API
func (s *UserService) GetUserProfile(userID string, sessionCookie *http.Cookie) (*models.UserProfile, error) {
	// Make API request using utils
	apiResponse, err := utils.MakeGETRequest(s.HTTPClient, s.BaseURL, "/users/profile/"+userID, sessionCookie)
	if err != nil {
		return nil, err
	}

	// Convert response to user profile
	var userProfile models.UserProfile
	if err := utils.ConvertAPIData(apiResponse.Data, &userProfile); err != nil {
		return nil, utils.NewGeneralError("Failed to parse user profile data", 500)
	}

	return &userProfile, nil
}

// GetUserPosts retrieves posts created by a specific user
func (s *UserService) GetUserPosts(userID string, limit, offset int, sortBy string, sessionCookie *http.Cookie) (*models.PaginatedPostsResponse, error) {
	// Build pagination parameters
	params := utils.BuildPaginationParams(limit, offset, sortBy)

	// Make API request using utils
	apiResponse, err := utils.MakeGETRequestWithParams(s.HTTPClient, s.BaseURL, "/users/posts/"+userID, params, sessionCookie)
	if err != nil {
		return nil, err
	}

	// Convert response to paginated posts
	var postsResponse models.PaginatedPostsResponse
	if err := utils.ConvertAPIData(apiResponse.Data, &postsResponse); err != nil {
		return nil, utils.NewGeneralError("Failed to parse posts data", 500)
	}

	return &postsResponse, nil
}

// GetUserLikedPosts retrieves posts liked by a specific user
func (s *UserService) GetUserLikedPosts(userID string, limit, offset int, sortBy string, sessionCookie *http.Cookie) (*models.PaginatedPostsResponse, error) {
	// Build pagination parameters
	params := utils.BuildPaginationParams(limit, offset, sortBy)

	// Make API request using utils
	apiResponse, err := utils.MakeGETRequestWithParams(s.HTTPClient, s.BaseURL, "/users/liked-posts/"+userID, params, sessionCookie)
	if err != nil {
		return nil, err
	}

	// Convert response to paginated posts
	var postsResponse models.PaginatedPostsResponse
	if err := utils.ConvertAPIData(apiResponse.Data, &postsResponse); err != nil {
		return nil, utils.NewGeneralError("Failed to parse liked posts data", 500)
	}

	return &postsResponse, nil
}

// GetUserCommentedPosts retrieves posts that a specific user has commented on
func (s *UserService) GetUserCommentedPosts(userID string, limit, offset int, sortBy string, sessionCookie *http.Cookie) (*models.PaginatedPostsResponse, error) {
	// Build pagination parameters
	params := utils.BuildPaginationParams(limit, offset, sortBy)

	// Make API request using utils
	apiResponse, err := utils.MakeGETRequestWithParams(s.HTTPClient, s.BaseURL, "/users/commented-posts/"+userID, params, sessionCookie)
	if err != nil {
		return nil, err
	}

	// Convert response to paginated posts
	var postsResponse models.PaginatedPostsResponse
	if err := utils.ConvertAPIData(apiResponse.Data, &postsResponse); err != nil {
		return nil, utils.NewGeneralError("Failed to parse commented posts data", 500)
	}

	return &postsResponse, nil
}
