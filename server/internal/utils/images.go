package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"real-time-forum/config"
	"real-time-forum/internal/models"
)

var allowedImageExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

// Checks file extension for allowed image types (case-insensitive)
func IsValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return allowedImageExt[ext]
}

func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// Returns os.File for writing
func CreateFile(path string) (*os.File, error) {
	// Make sure the directory exists
	err := os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		return nil, err
	}
	return os.Create(path)
}

// Wrapper around io.Copy for clarity
func CopyFile(dst *os.File, src multipart.File) (int64, error) {
	return io.Copy(dst, src)
}

func RemoveFileIfExists(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// validateImageFile validates the size and type of an uploaded image file
func validateImageFile(fileHeader *multipart.FileHeader) error {
	// Validate size
	if fileHeader.Size > 20*1024*1024 {
		return fmt.Errorf("file %s exceeds 20MB limit", fileHeader.Filename)
	}

	// Validate extension/type
	if !IsValidImageFile(fileHeader.Filename) {
		return fmt.Errorf("invalid file type: %s", fileHeader.Filename)
	}

	return nil
}

// generateUniqueImageFilename generates a unique filename for an image
func generateUniqueImageFilename(originalFilename string) (imageID, uniqueFilename, fullPath string) {
	imageID = GenerateUUIDToken()
	ext := GetFileExtension(originalFilename)
	uniqueFilename = fmt.Sprintf("%s%s", imageID, ext)
	fullPath = config.Config.UploadDir + uniqueFilename
	return
}

// saveImageToDisk saves the uploaded image file to disk
func saveImageToDisk(fileHeader *multipart.FileHeader, fullPath string) error {
	// Open source file
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	// Create destination file
	outFile, err := CreateFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}
	defer outFile.Close()

	// Copy data
	_, err = CopyFile(outFile, file)
	if err != nil {
		return fmt.Errorf("failed to save image data: %w", err)
	}

	return nil
}

// processAndSaveSingleImage validates, processes, and saves a single image file to disk
// Returns the PostImage metadata
func processAndSaveSingleImage(fileHeader *multipart.FileHeader) (models.PostImage, error) {
	// Validate the image file
	if err := validateImageFile(fileHeader); err != nil {
		return models.PostImage{}, err
	}

	// Generate unique filename
	imageID, uniqueFilename, fullPath := generateUniqueImageFilename(fileHeader.Filename)

	// Save image to disk
	if err := saveImageToDisk(fileHeader, fullPath); err != nil {
		return models.PostImage{}, err
	}

	return models.PostImage{
		ImageID:          imageID,
		ImageURL:         "/uploads/" + uniqueFilename,
		OriginalFilename: fileHeader.Filename,
	}, nil
}

// ProcessImageUploads handles validation and processing of uploaded images
func ProcessImageUploads(files []*multipart.FileHeader) ([]models.PostImage, error) {
	if len(files) == 0 {
		return []models.PostImage{}, nil
	}

	// Validate total number of images
	if len(files) > config.Config.MaxImagesPerPost {
		return nil, fmt.Errorf("maximum %d images allowed per post", config.Config.MaxImagesPerPost)
	}

	images := make([]models.PostImage, 0, len(files))

	for _, fileHeader := range files {
		image, err := processAndSaveSingleImage(fileHeader)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}
