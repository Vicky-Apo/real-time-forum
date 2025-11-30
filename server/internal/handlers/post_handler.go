package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"real-time-forum/config"
	"real-time-forum/internal/middleware"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/models"
	"real-time-forum/internal/utils"
)

// ...
// CRUD HANDLERS FOR POSTS
// ...

// CreatePostHandler creates a new post
func CreatePostHandler(pr *repository.PostsRepository, cr *repository.CategoryRepository, pir *repository.PostImagesRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// ---- Parse multipart form ----
		err := r.ParseMultipartForm(25 << 20) // 25MB max memory
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid multipart form")
			return
		}

		content := r.FormValue("content")
		if err := utils.ValidatePostContent(content); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse and validate categories
		categoryNames := r.MultipartForm.Value["categories"]
		categoryIDs, err := validateCategories(categoryNames, cr)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// ---- Parse and validate images ----
		images, err := utils.ProcessImageUploads(r.MultipartForm.File["images"])
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Create post - now returns lightweight response
		createResponse, err := pr.CreatePost(user.ID, content, categoryIDs, images)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create post")
			return
		}

		// Return lightweight response
		utils.RespondWithSuccess(w, http.StatusCreated, createResponse)
	}
}

// UpdatePostHandler updates an existing post
func UpdatePostHandler(pr *repository.PostsRepository, cr *repository.CategoryRepository, pir *repository.PostImagesRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := middleware.GetCurrentUser(r)

		postID := r.PathValue("id")
		if postID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Post ID is required")
			return
		}

		// Check ownership
		post, err := pr.GetPostByID(postID, user.ID)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Post not found")
			return
		}
		if post.UserID != user.ID {
			utils.RespondWithError(w, http.StatusForbidden, "You can only update your own posts")
			return
		}

		// Parse multipart form
		err = r.ParseMultipartForm(25 << 20) // 25MB
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid multipart form")
			return
		}

		// Parse new content/categories
		content := r.FormValue("content")
		if err := utils.ValidatePostContent(content); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		categoryNames := r.MultipartForm.Value["categories"]
		categoryIDs, err := validateCategories(categoryNames, cr)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// ---- Handle images ----

		// (A) Handle REMOVE images
		removeIDs := r.MultipartForm.Value["remove_image_ids[]"]
		removeImagesFromPost(removeIDs, pir)

		// (B) Handle ADD images
		if err = addImagesToPost(r.MultipartForm.File["images"], postID, pir); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Update the post's main content and categories
		err = pr.UpdatePost(postID, user.ID, content, categoryIDs)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update post")
			return
		}
		utils.RespondWithSuccess(w, http.StatusOK, nil)
	}
}

// DeletePostHandler deletes an existing post
func DeletePostHandler(pr *repository.PostsRepository, cr *repository.CategoryRepository, pir *repository.PostImagesRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// Extract post ID from URL path
		postID := r.PathValue("id")
		if postID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Post ID is required")
			return
		}

		// Pass userID to GetPostByID for ownership check
		post, err := pr.GetPostByID(postID, user.ID)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Post not found")
			return
		}
		// Check if the post belongs to the user
		if post.UserID != user.ID {
			utils.RespondWithError(w, http.StatusForbidden, "You can only delete your own posts")
			return
		}
		// Delete image files
		images, err := pir.DeleteAllImagesForPost(postID) // Deletes records and returns metadata
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete post images")
			return
		}
		deleteImageFilesFromDisk(images)
		// Delete the post
		err = pr.DeletePost(postID, user.ID)
		if err != nil {
			if err.Error() == "post not found" {
				utils.RespondWithError(w, http.StatusNotFound, "Post not found")
				return
			}
			if err.Error() == "unauthorized: you can only delete your own posts" {
				utils.RespondWithError(w, http.StatusForbidden, "You can only delete your own posts")
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete post")
			return
		}
		utils.RespondWithSuccess(w, http.StatusOK, nil)
	}
}

// GetAllPostsHandler retrieves all posts with pagination and sorting
func GetAllPostsHandler(pr *repository.PostsRepository) http.HandlerFunc {
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
func GetPostsByCategoryHandler(pr *repository.PostsRepository) http.HandlerFunc {
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
func GetSinglePostHandler(pr *repository.PostsRepository) http.HandlerFunc {
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

// GetUserPostsProfileHandler retrieves all posts by a specific user
func GetUserPostsProfileHandler(pr *repository.PostsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// Extract user ID from URL path
		userID := r.PathValue("id")
		if userID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "User ID is required")
			return
		}
		// Ensure user can only view their own posts
		if user.ID != userID {
			utils.RespondWithError(w, http.StatusForbidden, "You can only view your own posts")
			return
		}
		totalCount, err := pr.GetCountPostByUser(userID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve posts count")
			return
		}
		limit, offset := utils.ParsePaginationParams(r)
		sortOptions := utils.ParsePostSortOptions(r)
		//Pass both targetUserID and currentUserID to repository
		posts, err := pr.GetPostsByUser(userID, limit, offset, user.ID, sortOptions)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user posts")
			return
		}
		utils.RespondWithPaginatedPosts(w, posts, totalCount, limit, offset)
	}
}

// GetUserLikedPostsProfileHandler retrieves all posts liked by a specific user
func GetUserLikedPostsProfileHandler(pr *repository.PostsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// Extract user ID from URL path
		userID := r.PathValue("id")
		if userID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "User ID is required")
			return
		}
		// Ensure user can only view their own liked posts
		if user.ID != userID {
			utils.RespondWithError(w, http.StatusForbidden, "You can only view your own liked posts")
			return
		}
		// Parse pagination parameters
		limit, offset := utils.ParsePaginationParams(r)
		// sorting options
		sortOptions := utils.ParsePostSortOptions(r)
		totalCount, err := pr.GetCountLikedPostByUser(userID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve liked posts count")
			return
		}
		// Pass both targetUserID and currentUserID to repository
		posts, err := pr.GetPostsLikedByUser(userID, limit, offset, user.ID, sortOptions)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve liked posts")
			return
		}
		utils.RespondWithPaginatedPosts(w, posts, totalCount, limit, offset)
	}
}

// GetUserCommentedPostsProfileHandler retrieves all posts that a specific user has commented on
func GetUserCommentedPostsProfileHandler(pr *repository.PostsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

		// Extract user ID from URL path
		userID := r.PathValue("id")
		if userID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "User ID is required")
			return
		}

		// Ensure user can only view their own commented posts
		if user.ID != userID {
			utils.RespondWithError(w, http.StatusForbidden, "You can only view your own commented posts")
			return
		}

		// Parse pagination parameters
		limit, offset := utils.ParsePaginationParams(r)
		// Parse sort options from query parameters
		sortOptions := utils.ParsePostSortOptions(r)
		// Get total count of posts user has commented on
		totalCount, err := pr.GetCountCommentedPostByUser(userID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve commented posts count")
			return
		}

		// Get posts that user has commented on
		posts, err := pr.GetPostsCommentedByUser(userID, limit, offset, user.ID, sortOptions)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve commented posts")
			return
		}

		// Respond with paginated posts
		utils.RespondWithPaginatedPosts(w, posts, totalCount, limit, offset)
	}
}

// ...
// HELPER FUNCTIONS
// ...

// validateCategories validates category names and returns their IDs
func validateCategories(categoryNames []string, cr *repository.CategoryRepository) ([]string, error) {
	if len(categoryNames) == 0 {
		return nil, fmt.Errorf("at least one category is required")
	}

	var categoryIDs []string
	for _, categoryName := range categoryNames {
		categoryID, err := cr.GetCategoryID(categoryName)
		if err != nil {
			return nil, fmt.Errorf("invalid category: %s", categoryName)
		}
		categoryIDs = append(categoryIDs, categoryID)
	}

	if len(categoryIDs) < config.Config.MinCategories {
		return nil, fmt.Errorf("minimum %d category required", config.Config.MinCategories)
	}

	if len(categoryIDs) > config.Config.MaxCategories {
		return nil, fmt.Errorf("maximum %d categories allowed", config.Config.MaxCategories)
	}

	return categoryIDs, nil
}

// deleteImageFilesFromDisk deletes image files from disk
func deleteImageFilesFromDisk(images []models.PostImage) {
	for _, img := range images {
		filePath := "." + img.ImageURL // assuming "/uploads/xyz" -> "./uploads/xyz"
		err := utils.RemoveFileIfExists(filePath)
		if err != nil {
			// Log and continue
			fmt.Printf("Warning: failed to delete image file %s: %v\n", filePath, err)
		}
	}
}

// removeImagesFromPost removes images from a post by their IDs
func removeImagesFromPost(removeIDs []string, pir *repository.PostImagesRepository) {
	for _, imgID := range removeIDs {
		// Delete from DB and get info for file deletion
		img, err := pir.DeleteImageByID(imgID)
		if err != nil {
			// Log, but don't fail whole update
			fmt.Printf("Failed to delete image %s: %v\n", imgID, err)
			continue
		}
		// Remove file from disk
		deleteImageFilesFromDisk([]models.PostImage{*img})
	}
}

// addImagesToPost processes and saves images for an existing post
func addImagesToPost(files []*multipart.FileHeader, postID string, pir *repository.PostImagesRepository) error {
	if len(files) == 0 {
		return nil
	}

	if len(files) > config.Config.MaxImagesPerPost {
		return fmt.Errorf("maximum %d images allowed per post", config.Config.MaxImagesPerPost)
	}

	for _, fileHeader := range files {
		// Process and save the image file
		images, err := utils.ProcessImageUploads([]*multipart.FileHeader{fileHeader})
		if err != nil {
			return err
		}

		// Insert image record to database
		image := images[0]
		if err = pir.SaveImageRecord(postID, image.ImageID, image.ImageURL, image.OriginalFilename); err != nil {
			return fmt.Errorf("failed to save image record: %w", err)
		}
	}

	return nil
}
