package handlers

import (
	"log"
	"net/http"

	"frontend-service/config"
	"frontend-service/internal/services"
)

type OAuthHandler struct {
	authService     *services.AuthService
	oauthService    *services.OAuthService
	templateService *services.TemplateService
	config          *config.Config
	errorHandler    *SimpleErrorHandler
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(authService *services.AuthService, oauthService *services.OAuthService, templateService *services.TemplateService, cfg *config.Config) *OAuthHandler {
	return &OAuthHandler{
		authService:     authService,
		oauthService:    oauthService,
		templateService: templateService,
		config:          cfg,
		errorHandler:    NewSimpleErrorHandler(templateService),
	}
}

// ================================
// OAUTH INITIATION HANDLERS
// ================================

// Redirects user to backend GitHub OAuth flow
func (h *OAuthHandler) ServeGitHubLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler.ShowError(w, "Method Not Allowed", "This method is not allowed for GitHub login.")
		return
	}

	// Get the GitHub OAuth URL from backend
	githubURL := h.oauthService.GetGitHubLoginURL()

	log.Printf("Redirecting user to GitHub OAuth: %s", githubURL)

	// Redirect user to backend OAuth endpoint
	http.Redirect(w, r, githubURL, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) ServeGoogleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler.ShowError(w, "Method Not Allowed", "This method is not allowed for Google login.")
		return
	}

	// Get the Google OAuth URL from backend
	googleURL := h.oauthService.GetGoogleLoginURL()

	log.Printf("Redirecting user to Google OAuth: %s", googleURL)

	// Redirect user to backend OAuth endpoint
	http.Redirect(w, r, googleURL, http.StatusTemporaryRedirect)
}
