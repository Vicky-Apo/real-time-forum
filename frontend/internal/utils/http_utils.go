// internal/utils/http_utils.go
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"frontend-service/internal/models"
)

// MakeGETRequest - Simple GET request with session cookie
func MakeGETRequest(httpClient *http.Client, baseURL, endpoint string, sessionCookie *http.Cookie) (*models.APIResponse, error) {
	url := baseURL + endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}

	return executeRequest(httpClient, req)
}

// MakeGETRequestWithParams - GET request with query parameters
func MakeGETRequestWithParams(httpClient *http.Client, baseURL, endpoint string, params map[string]string, sessionCookie *http.Cookie) (*models.APIResponse, error) {
	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	urlParams := url.Values{}
	for key, value := range params {
		if value != "" {
			urlParams.Add(key, value)
		}
	}
	u.RawQuery = urlParams.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}

	return executeRequest(httpClient, req)
}

// MakePOSTRequest - POST request with JSON data
func MakePOSTRequest(httpClient *http.Client, baseURL, endpoint string, data interface{}, sessionCookie *http.Cookie) (*models.APIResponse, error) {
	return makeJSONRequest(httpClient, "POST", baseURL, endpoint, data, sessionCookie)
}

// MakePUTRequest - PUT request with JSON data
func MakePUTRequest(httpClient *http.Client, baseURL, endpoint string, data interface{}, sessionCookie *http.Cookie) (*models.APIResponse, error) {
	return makeJSONRequest(httpClient, "PUT", baseURL, endpoint, data, sessionCookie)
}

// MakeDELETERequest - DELETE request
func MakeDELETERequest(httpClient *http.Client, baseURL, endpoint string, sessionCookie *http.Cookie) (*models.APIResponse, error) {
	url := baseURL + endpoint

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}

	return executeRequest(httpClient, req)
}

// ForwardMultipartRequest - Forward raw request (for file uploads)
func ForwardMultipartRequest(httpClient *http.Client, baseURL, endpoint, method string, sourceRequest *http.Request, sessionCookie *http.Cookie) (*models.APIResponse, error) {
	url := baseURL + endpoint

	req, err := http.NewRequest(method, url, sourceRequest.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", sourceRequest.Header.Get("Content-Type"))

	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}

	return executeRequest(httpClient, req)
}

// makeJSONRequest - Internal helper for POST/PUT with JSON
func makeJSONRequest(httpClient *http.Client, method, baseURL, endpoint string, data interface{}, sessionCookie *http.Cookie) (*models.APIResponse, error) {
	url := baseURL + endpoint

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}

	return executeRequest(httpClient, req)
}

// executeRequest - Internal helper that executes request and handles response
func executeRequest(httpClient *http.Client, req *http.Request) (*models.APIResponse, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle HTTP status codes
	if err := HandleHTTPStatus(resp.StatusCode, body); err != nil {
		return nil, err
	}

	var apiResponse models.APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("API error: %s", apiResponse.Error)
	}

	return &apiResponse, nil
}

// BuildPaginationParams - Create pagination parameters map
func BuildPaginationParams(limit, offset int, sortBy string) map[string]string {
	params := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
	}

	if sortBy != "" {
		params["sort"] = sortBy
	}

	return params
}

// ConvertAPIData - Convert API response data to target struct
func ConvertAPIData(apiData interface{}, target interface{}) error {
	dataBytes, err := json.Marshal(apiData)
	if err != nil {
		return fmt.Errorf("failed to marshal API data: %w", err)
	}

	if err := json.Unmarshal(dataBytes, target); err != nil {
		return fmt.Errorf("failed to unmarshal to target struct: %w", err)
	}

	return nil
}
