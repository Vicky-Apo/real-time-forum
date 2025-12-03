// internal/services/auth_service.go
package services

import (
	"net/http"

	"frontend-service/config"
	"frontend-service/internal/models"
	"frontend-service/internal/utils"
	"frontend-service/internal/validations"
)

type AuthService struct {
	*BaseClient
	SessionName string // Store session cookie name from config
}

// NewAuthService creates a new auth service
func NewAuthService(baseClient *BaseClient, cfg *config.Config) *AuthService {
	return &AuthService{
		BaseClient:  baseClient,
		SessionName: cfg.SessionName,
	}
}

// GetSessionName returns the session cookie name (implements AuthServiceInterface)
func (s *AuthService) GetSessionName() string {
	return s.SessionName
}

// RegisterUser registers a new user via the backend API
func (s *AuthService) RegisterUser(formData models.UserRegistration) error {
	// Frontend validation
	if err := validations.ValidateUserInput(formData.Username, formData.Email, formData.Password); err != nil {
		return utils.NewGeneralError("Invalid user input: "+err.Error(), 400)
	}

	// Make API request using utils
	_, err := utils.MakePOSTRequest(s.HTTPClient, s.BaseURL, "/auth/register", formData, nil)
	return err
}

// LoginUser logs in a user via the backend API
func (s *AuthService) LoginUser(formData models.UserLogin) (*models.User, string, error) {
	// Frontend validation
	if err := validations.ValidateEmail(formData.Email); err != nil {
		return nil, "", utils.NewGeneralError("Invalid email format: "+err.Error(), 400)
	}
	if err := validations.ValidatePassword(formData.Password); err != nil {
		return nil, "", utils.NewGeneralError("Invalid password format: "+err.Error(), 400)
	}

	// Make API request using utils
	apiResponse, err := utils.MakePOSTRequest(s.HTTPClient, s.BaseURL, "/auth/login", formData, nil)
	if err != nil {
		return nil, "", err
	}

	// Convert response to login data
	var loginResponse struct {
		User      models.User `json:"user"`
		SessionID string      `json:"session_id"`
	}

	if err := utils.ConvertAPIData(apiResponse.Data, &loginResponse); err != nil {
		return nil, "", utils.NewGeneralError("Failed to parse login response", 500)
	}

	return &loginResponse.User, loginResponse.SessionID, nil
}

// LogoutUser logs out a user via the backend API
func (s *AuthService) LogoutUser(sessionID string) error {
	// Create session cookie for the request
	sessionCookie := &http.Cookie{
		Name:  s.SessionName,
		Value: sessionID,
	}

	// Make API request using utils
	_, err := utils.MakePOSTRequest(s.HTTPClient, s.BaseURL, "/auth/logout", nil, sessionCookie)
	return err
}

// ValidateSession validates a session ID with the backend API (implements AuthServiceInterface)
func (s *AuthService) ValidateSession(sessionID string) (*models.User, error) {
	// Create session cookie for the request
	sessionCookie := &http.Cookie{
		Name:  s.SessionName,
		Value: sessionID,
	}

	// Make API request using utils
	apiResponse, err := utils.MakePOSTRequest(s.HTTPClient, s.BaseURL, "/auth/me", nil, sessionCookie)
	if err != nil {
		return nil, err
	}

	// Convert response to user data
	var user models.User
	if err := utils.ConvertAPIData(apiResponse.Data, &user); err != nil {
		return nil, utils.NewGeneralError("Failed to parse user data", 500)
	}

	return &user, nil
}
