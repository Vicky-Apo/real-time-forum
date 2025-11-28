package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
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
