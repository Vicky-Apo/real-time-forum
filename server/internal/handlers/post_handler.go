package handlers

import (
	"fmt"
	"net/http"

	"platform.zone01.gr/git/gpapadopoulos/forum/config"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/middleware"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/models"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/repository"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/utils"
)

// ...
// CRUD HANDLERS FOR POSTS
// ...
// REPLACE the CreatePostHandler in post_handler.go:

func CreatePostHandler(pr *repository.PostsRepository, cr *repository.CategoryRepository, pir *repository.PostImagesRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)
		if user == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// ---- Parse multipart form ----
		err := r.ParseMultipartForm(25 << 20) // 25MB max memory
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid multipart form")
			return
		}

		content := r.FormValue("content")
		if utils.ValidatePostContent(content) != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid post content")
			return
		}

		// Parse and validate categories
		categoryNames := r.MultipartForm.Value["categories"]
		if len(categoryNames) == 0 {
			utils.RespondWithError(w, http.StatusBadRequest, "At least one category is required")
			return
		}

		// Validate and get category IDs
		var categoryIDs []string
		for _, categoryName := range categoryNames {
			categoryID, err := cr.GetCategoryID(categoryName)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid category: "+categoryName)
				return
			}
			categoryIDs = append(categoryIDs, categoryID)
		}
		if len(categoryIDs) < config.Config.MinCategories {
			utils.RespondWithError(w, http.StatusBadRequest,
				fmt.Sprintf("Minimum %d category required", config.Config.MinCategories))
			return
		}

		if len(categoryIDs) > config.Config.MaxCategories {
			utils.RespondWithError(w, http.StatusBadRequest,
				fmt.Sprintf("Maximum %d categories allowed", config.Config.MaxCategories))
			return
		}

		// ---- Parse and validate images ----
		var images []models.PostImage
		files := r.MultipartForm.File["images"]
		if len(files) > 0 {
			if len(files) > config.Config.MaxImagesPerPost {
				utils.RespondWithError(w, http.StatusBadRequest,
					fmt.Sprintf("Maximum %d images allowed per post", config.Config.MaxImagesPerPost))
				return
			}
			for _, fileHeader := range files {
				// Validate size
				if fileHeader.Size > 20*1024*1024 {
					utils.RespondWithError(w, http.StatusBadRequest,
						fmt.Sprintf("File %s exceeds 20MB limit", fileHeader.Filename))
					return
				}

				// Validate extension/type
				valid := utils.IsValidImageFile(fileHeader.Filename)
				if !valid {
					utils.RespondWithError(w, http.StatusBadRequest,
						fmt.Sprintf("Invalid file type: %s", fileHeader.Filename))
					return
				}

				// Open file
				file, err := fileHeader.Open()
				if err != nil {
					utils.RespondWithError(w, http.StatusInternalServerError, "Failed to open image file")
					return
				}
				defer file.Close()

				// Generate unique filename and save to disk
				imageID := utils.GenerateUUIDToken()
				ext := utils.GetFileExtension(fileHeader.Filename)
				uniqueFilename := fmt.Sprintf("%s%s", imageID, ext)
				uploadPath := config.Config.UploadDir // e.g., "./uploads/"
				fullPath := uploadPath + uniqueFilename

				outFile, err := utils.CreateFile(fullPath) // wraps os.Create with error checks
				if err != nil {
					utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save image")
					return
				}
				defer outFile.Close()

				_, err = utils.CopyFile(outFile, file) // wraps io.Copy
				if err != nil {
					utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save image data")
					return
				}

				images = append(images, models.PostImage{
					ImageID:          imageID,
					ImageURL:         "/uploads/" + uniqueFilename, // adjust if you serve differently
					OriginalFilename: fileHeader.Filename,
				})
			}
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
		if user == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}
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
		if utils.ValidatePostContent(content) != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid post content")
			return
		}
		categoryNames := r.MultipartForm.Value["categories"]
		if len(categoryNames) == 0 {
			utils.RespondWithError(w, http.StatusBadRequest, "At least one category is required")
			return
		}

		var categoryIDs []string
		for _, categoryName := range categoryNames {
			categoryID, err := cr.GetCategoryID(categoryName)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid category: "+categoryName)
				return
			}
			categoryIDs = append(categoryIDs, categoryID)
		}
		if len(categoryIDs) < config.Config.MinCategories {
			utils.RespondWithError(w, http.StatusBadRequest,
				fmt.Sprintf("Minimum %d category required", config.Config.MinCategories))
			return
		}
		if len(categoryIDs) > config.Config.MaxCategories {
			utils.RespondWithError(w, http.StatusBadRequest,
				fmt.Sprintf("Maximum %d categories allowed", config.Config.MaxCategories))
			return
		}

		// ---- Handle images ----

		// (A) Handle REMOVE images
		removeIDs := r.MultipartForm.Value["remove_image_ids[]"]
		for _, imgID := range removeIDs {
			// Delete from DB and get info for file deletion
			img, err := pir.DeleteImageByID(imgID)
			if err != nil {
				// Log, but don't fail whole update
				fmt.Printf("Failed to delete image %s: %v\n", imgID, err)
				continue
			}
			// Remove file from disk
			filePath := "." + img.ImageURL
			if err := utils.RemoveFileIfExists(filePath); err != nil {
				fmt.Printf("Failed to delete image file %s: %v\n", filePath, err)
			}
		}

		// (B) Handle ADD images (same logic as Create)
		files := r.MultipartForm.File["images"]
		if len(files) > 0 {
			if len(files) > config.Config.MaxImagesPerPost {
				utils.RespondWithError(w, http.StatusBadRequest,
					fmt.Sprintf("Maximum %d images allowed per post", config.Config.MaxImagesPerPost))
				return
			}
			for _, fileHeader := range files {
				if fileHeader.Size > 20*1024*1024 {
					utils.RespondWithError(w, http.StatusBadRequest,
						fmt.Sprintf("File %s exceeds 20MB limit", fileHeader.Filename))
					return
				}
				valid := utils.IsValidImageFile(fileHeader.Filename)
				if !valid {
					utils.RespondWithError(w, http.StatusBadRequest,
						fmt.Sprintf("Invalid file type: %s", fileHeader.Filename))
					return
				}
				file, err := fileHeader.Open()
				if err != nil {
					utils.RespondWithError(w, http.StatusInternalServerError, "Failed to open image file")
					return
				}
				defer file.Close()
				imageID := utils.GenerateUUIDToken()
				ext := utils.GetFileExtension(fileHeader.Filename)
				uniqueFilename := fmt.Sprintf("%s%s", imageID, ext)
				uploadPath := config.Config.UploadDir
				fullPath := uploadPath + uniqueFilename
				outFile, err := utils.CreateFile(fullPath)
				if err != nil {
					utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save image")
					return
				}
				defer outFile.Close()
				_, err = utils.CopyFile(outFile, file)
				if err != nil {
					utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save image data")
					return
				}
				// Insert image record
				err = pir.SaveImageRecord(postID, imageID, "/uploads/"+uniqueFilename, fileHeader.Filename)
				if err != nil {
					utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save image record")
					return
				}
			}
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
		if user == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}
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
		// --- NEW: Delete image files ---
		images, err := pir.DeleteAllImagesForPost(postID) // Deletes records and returns metadata
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete post images")
			return
		}

		for _, img := range images {
			filePath := "." + img.ImageURL // assuming "/uploads/xyz" -> "./uploads/xyz"
			err := utils.RemoveFileIfExists(filePath)
			if err != nil {
				// Log and continue
				fmt.Printf("Warning: failed to delete image file %s: %v\n", filePath, err)
			}
		}
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
