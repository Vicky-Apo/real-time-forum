// internal/services/post_service.go
package services

import (
	"net/http"

	"frontend-service/internal/models"
	"frontend-service/internal/utils"
)

type PostService struct {
	*BaseClient
}

// NewPostService creates a new post service
func NewPostService(baseClient *BaseClient) *PostService {
	return &PostService{
		BaseClient: baseClient,
	}
}

// GetAllPosts retrieves posts from the backend API with pagination and sorting
func (s *PostService) GetAllPosts(limit, offset int, sortBy string, sessionCookie *http.Cookie) (*models.PaginatedPostsResponse, error) {
	// Build pagination parameters
	params := utils.BuildPaginationParams(limit, offset, sortBy)

	// Make API request using utils
	apiResponse, err := utils.MakeGETRequestWithParams(s.HTTPClient, s.BaseURL, "/posts", params, sessionCookie)
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

// GetPostsByCategory retrieves posts filtered by category with pagination and sorting
func (s *PostService) GetPostsByCategory(categoryID string, limit, offset int, sortBy string, sessionCookie *http.Cookie) (*models.PaginatedPostsResponse, error) {
	// Build pagination parameters
	params := utils.BuildPaginationParams(limit, offset, sortBy)

	// Make API request using utils
	apiResponse, err := utils.MakeGETRequestWithParams(s.HTTPClient, s.BaseURL, "/posts/by-category/"+categoryID, params, sessionCookie)
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

// GetSinglePostWithComments retrieves a single post and its comments with pagination
func (s *PostService) GetSinglePostWithComments(postID string, limit, offset int, sortBy string, sessionCookie *http.Cookie) (*models.Post, *models.PaginatedCommentsResponse, error) {
	// Get the post
	post, err := s.GetSinglePost(postID, sessionCookie)
	if err != nil {
		return nil, nil, err
	}

	// Get the comments with pagination
	commentsResponse, err := s.GetPostComments(postID, limit, offset, sortBy, sessionCookie)
	if err != nil {
		return nil, nil, err
	}

	return post, commentsResponse, nil
}

// GetSinglePost retrieves a single post by ID
func (s *PostService) GetSinglePost(postID string, sessionCookie *http.Cookie) (*models.Post, error) {
	// Make API request using utils
	apiResponse, err := utils.MakeGETRequest(s.HTTPClient, s.BaseURL, "/posts/view/"+postID, sessionCookie)
	if err != nil {
		return nil, err
	}

	// Convert response to post
	var post models.Post
	if err := utils.ConvertAPIData(apiResponse.Data, &post); err != nil {
		return nil, utils.NewGeneralError("Failed to parse post data", 500)
	}

	return &post, nil
}

// GetPostComments retrieves comments for a post with pagination and sorting
func (s *PostService) GetPostComments(postID string, limit, offset int, sortBy string, sessionCookie *http.Cookie) (*models.PaginatedCommentsResponse, error) {
	// Build pagination parameters
	params := utils.BuildPaginationParams(limit, offset, sortBy)

	// Make API request using utils
	apiResponse, err := utils.MakeGETRequestWithParams(s.HTTPClient, s.BaseURL, "/comments/for-post/"+postID, params, sessionCookie)
	if err != nil {
		return nil, err
	}

	// Convert response to paginated comments
	var commentsResponse models.PaginatedCommentsResponse
	if err := utils.ConvertAPIData(apiResponse.Data, &commentsResponse); err != nil {
		return nil, utils.NewGeneralError("Failed to parse comments data", 500)
	}

	return &commentsResponse, nil
}

// CreatePost creates a new post by forwarding multipart form data
func (s *PostService) CreatePost(r *http.Request, sessionCookie *http.Cookie) (string, error) {
	// Forward the multipart request using utils
	apiResponse, err := utils.ForwardMultipartRequest(s.HTTPClient, s.BaseURL, "/posts/create", "POST", r, sessionCookie)
	if err != nil {
		return "", err
	}

	// Extract post ID from response
	var createResponse struct {
		PostID string `json:"post_id"`
	}

	if err := utils.ConvertAPIData(apiResponse.Data, &createResponse); err != nil {
		return "", utils.NewGeneralError("Failed to parse create response", 500)
	}

	return createResponse.PostID, nil
}

// UpdatePost updates an existing post by forwarding multipart form data
func (s *PostService) UpdatePost(r *http.Request, postID string, sessionCookie *http.Cookie) error {
	// Forward the multipart request using utils
	_, err := utils.ForwardMultipartRequest(s.HTTPClient, s.BaseURL, "/posts/edit/"+postID, "PUT", r, sessionCookie)
	return err
}

// DeletePost deletes a post by ID
func (s *PostService) DeletePost(postID string, sessionCookie *http.Cookie) error {
	// Make API request using utils
	_, err := utils.MakeDELETERequest(s.HTTPClient, s.BaseURL, "/posts/remove/"+postID, sessionCookie)
	return err
}

// ForwardCreatePost - Legacy method for backward compatibility
// DEPRECATED: Use CreatePost instead
func (s *PostService) ForwardCreatePost(r *http.Request, sessionCookie *http.Cookie) (postID string, apiError string, status int) {
	postID, err := s.CreatePost(r, sessionCookie)
	if err != nil {
		// Convert new error format to old format for backward compatibility
		if appErr, ok := err.(*utils.AppError); ok {
			return "", appErr.Message, appErr.Code
		}
		return "", err.Error(), 500
	}
	return postID, "", 201
}
