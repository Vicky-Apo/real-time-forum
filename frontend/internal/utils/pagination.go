package utils

import (
	"net/http"
	"strconv"
	"strings"
)

// Pagination constants for frontend
const (
	DefaultPageSize = 20
	MaxPageSize     = 50
	MinPageSize     = 1
)

// ParsePaginationParams extracts and validates pagination parameters from HTTP request
func ParsePaginationParams(r *http.Request) (limit, offset int) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Convert to integers
	limit, _ = strconv.Atoi(limitStr)
	offset, _ = strconv.Atoi(offsetStr)

	// Validate and normalize using the model helper
	return ValidatePaginationParams(limit, offset)
}

// ValidatePaginationParams validates and normalizes pagination parameters
func ValidatePaginationParams(limit, offset int) (int, int) {
	// Set defaults and validate limit using constants
	if limit <= 0 {
		limit = DefaultPageSize
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}

	// Validate offset
	if offset < 0 {
		offset = 0
	}

	return limit, offset
}

// ParseSortFromRequest extracts sort parameter from HTTP request
func ParseSortFromRequest(r *http.Request, defaultSort string) string {
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		return defaultSort
	}
	return sortBy
}

// BuildPaginationURL creates a URL with pagination parameters
func BuildPaginationURL(baseURL string, limit, offset int) string {
	if baseURL == "" {
		baseURL = "?"
	} else if baseURL[len(baseURL)-1] != '?' && baseURL[len(baseURL)-1] != '&' {
		if strings.Contains(baseURL, "?") {
			baseURL += "&"
		} else {
			baseURL += "?"
		}
	}

	return baseURL + "limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)
}

// CalculateOffset calculates offset from page number
func CalculateOffset(page, limit int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * limit
}
