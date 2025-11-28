// internal/utils/validation_utils.go
package utils

import (
	"fmt"
	"net/http"
)

// ValidateHTTPMethod - Check if HTTP method is allowed
func ValidateHTTPMethod(r *http.Request, allowedMethods []string) error {
	for _, method := range allowedMethods {
		if r.Method == method {
			return nil
		}
	}
	return NewGeneralError("Method not allowed", http.StatusMethodNotAllowed)
}

// ValidateResourceID - Check if resource ID is not empty
func ValidateResourceID(id, resourceName string) error {
	if id == "" {
		return NewHomeError(fmt.Sprintf("%s ID is required", resourceName))
	}
	return nil
}

// ValidateFormData - Parse and validate form data
func ValidateFormData(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return NewGeneralError("Invalid form data provided", http.StatusBadRequest)
	}
	return nil
}

// ValidateSortParameter - Check if sort parameter is valid
func ValidateSortParameter(sortBy string, validSorts []string) error {
	if sortBy == "" {
		return nil // Default will be used
	}

	for _, valid := range validSorts {
		if sortBy == valid {
			return nil
		}
	}

	return NewGeneralError(
		fmt.Sprintf("Invalid sort option. Valid options: %v", validSorts),
		http.StatusBadRequest,
	)
}

// Common validation constants
var (
	ValidPostSorts    = []string{"newest", "oldest", "likes", "comments"}
	ValidCommentSorts = []string{"newest", "oldest", "likes"}
	GETMethod         = []string{http.MethodGet}
	POSTMethod        = []string{http.MethodPost}
	GETAndPOST        = []string{http.MethodGet, http.MethodPost}
)
