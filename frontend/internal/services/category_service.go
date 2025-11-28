// internal/services/category_service.go
package services

import (
	"frontend-service/internal/models"
	"frontend-service/internal/utils"
)

type CategoryService struct {
	*BaseClient
}

// NewCategoryService creates a new category service
func NewCategoryService(baseClient *BaseClient) *CategoryService {
	return &CategoryService{
		BaseClient: baseClient,
	}
}

// GetCategories retrieves categories from the backend API
func (s *CategoryService) GetCategories() ([]models.Category, error) {
	// Make API request using utils
	apiResponse, err := utils.MakeGETRequest(s.HTTPClient, s.BaseURL, "/categories", nil)
	if err != nil {
		return nil, err
	}

	// Convert response to categories slice
	var categories []models.Category
	if err := utils.ConvertAPIData(apiResponse.Data, &categories); err != nil {
		return nil, utils.NewGeneralError("Failed to parse categories data", 500)
	}

	return categories, nil
}
