package services

import (
	"frontend-service/config"
)

type OAuthService struct {
	*BaseClient
	config *config.Config
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(baseClient *BaseClient, cfg *config.Config) *OAuthService {
	return &OAuthService{
		BaseClient: baseClient,
		config:     cfg,
	}
}

// ================================
// OAUTH FLOW METHODS
// ================================

// GetGitHubLoginURL returns the URL to start GitHub OAuth flow
// This should point to the backend's web OAuth endpoint
func (s *OAuthService) GetGitHubLoginURL() string {
	// BaseURL already includes /api, so just add the endpoint path
	return s.BaseURL + "/auth/github/login"

}

func (s *OAuthService) GetGoogleLoginURL() string {
	return s.BaseURL + "/auth/google/login"
}
