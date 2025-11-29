package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"real-time-forum/config"
	"real-time-forum/internal/models"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/utils"
)

type OAuthHandler struct {
	oauthRepo   *repository.OAuthRepository
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	config      *config.AppConfig
}

func NewOAuthHandler(oauthRepo *repository.OAuthRepository, userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository, cfg *config.AppConfig) *OAuthHandler {
	return &OAuthHandler{
		oauthRepo:   oauthRepo,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		config:      cfg,
	}
}

// ================================
// GITHUB OAUTH FLOW HANDLERS (WEB ONLY)
// ================================

// ServeGitHubLogin initiates GitHub OAuth flow (WEB ONLY)
func (h *OAuthHandler) ServeGitHubLogin(w http.ResponseWriter, r *http.Request) {

	// Create OAuth state for CSRF protection
	state, err := h.oauthRepo.CreateOAuthState("github")
	if err != nil {
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=state_failed", http.StatusSeeOther)
		return
	}

	// Build GitHub authorization URL (WEB CALLBACK ONLY)
	authURL := h.buildGitHubAuthURL(state.StateID)

	// Redirect user to GitHub
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// ServeGitHubCallback handles GitHub OAuth callback (WEB ONLY)
func (h *OAuthHandler) ServeGitHubCallback(w http.ResponseWriter, r *http.Request) {

	// Get parameters from URL
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errorParam := r.URL.Query().Get("error")

	// Handle user denial
	if errorParam == "access_denied" {
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=github_cancelled", http.StatusSeeOther)
		return
	}

	// Validate required parameters
	if code == "" || state == "" {
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=missing_parameters", http.StatusSeeOther)
		return
	}

	// Validate state (CSRF protection)
	err := h.oauthRepo.ValidateOAuthState(state, "github")
	if err != nil {
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=invalid_state", http.StatusSeeOther)
		return
	}

	// Exchange authorization code for access token
	accessToken, err := h.exchangeGitHubCode(code)
	if err != nil {
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=token_exchange_failed", http.StatusSeeOther)
		return
	}

	// Get user information from GitHub
	gitHubUser, err := h.getGitHubUserInfo(accessToken)
	if err != nil {
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=user_info_failed", http.StatusSeeOther)
		return
	}

	// Process the OAuth authentication
	result, err := h.processGitHubAuthentication(gitHubUser, accessToken)
	if err != nil {
		// Check if it's an email conflict error
		if strings.Contains(err.Error(), "email_conflict:") {
			http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=email_conflict", http.StatusSeeOther)
			return
		}
		// Other authentication errors
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=auth_failed", http.StatusSeeOther)
		return
	}

	// Create session for the authenticated user
	session, err := h.sessionRepo.CreateSession(result.User.ID, r.RemoteAddr)
	if err != nil {
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=session_failed", http.StatusSeeOther)
		return
	}

	// Set session cookie
	utils.SetSessionCookie(session.SessionID, w, r, session.ExpiresAt)

	// Redirect to frontend with success parameters
	var redirectURL string
	if result.IsNewUser {
		redirectURL = h.config.FrontendBaseURL + "/?welcome=github"
	} else if result.IsLinkedAccount {
		redirectURL = h.config.FrontendBaseURL + "/?linked=github"
	} else {
		redirectURL = h.config.FrontendBaseURL + "/"
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// ================================
// GITHUB API FUNCTIONS
// ================================

// buildGitHubAuthURL builds the GitHub authorization URL
func (h *OAuthHandler) buildGitHubAuthURL(state string) string {
	params := url.Values{}
	params.Add("client_id", config.Config.GitHubClientID)
	params.Add("redirect_uri", h.config.BackendBaseURL+"/api/auth/github/callback")
	params.Add("scope", "user:email")
	params.Add("state", state)

	return "https://github.com/login/oauth/authorize?" + params.Encode()
}

// exchangeGitHubCode exchanges authorization code for access token
func (h *OAuthHandler) exchangeGitHubCode(code string) (string, error) {
	// Prepare form data
	data := url.Values{}
	data.Set("client_id", config.Config.GitHubClientID)
	data.Set("client_secret", config.Config.GitHubClientSecret)
	data.Set("code", code)

	// Create request
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("token exchange request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %w", err)
	}

	// Parse JSON response
	var tokenResp models.GitHubTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	// Check for errors
	if tokenResp.Error != "" {
		return "", fmt.Errorf("GitHub OAuth error: %s - %s", tokenResp.Error, tokenResp.ErrorDescription)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("no access token received from GitHub")
	}

	return tokenResp.AccessToken, nil
}

// getGitHubUserInfo fetches user information from GitHub API
func (h *OAuthHandler) getGitHubUserInfo(accessToken string) (*models.GitHubUser, error) {
	// Create request to GitHub user API
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user info request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	var gitHubUser models.GitHubUser
	if err := json.Unmarshal(body, &gitHubUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// Get user's primary email if not included
	if gitHubUser.Email == "" {
		email, err := h.getGitHubUserEmail(accessToken)
		if err == nil {
			gitHubUser.Email = email
		}
	}

	return &gitHubUser, nil
}

// getGitHubUserEmail fetches user's primary email from GitHub API
func (h *OAuthHandler) getGitHubUserEmail(accessToken string) (string, error) {
	// Create request to GitHub emails API
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create emails request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("emails request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read emails response: %w", err)
	}

	var emails []models.GitHubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", fmt.Errorf("failed to parse emails: %w", err)
	}

	// Find primary verified email
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no primary verified email found")
}

// ================================
// AUTHENTICATION PROCESSING
// ================================

// processGitHubAuthentication handles the core OAuth authentication logic
func (h *OAuthHandler) processGitHubAuthentication(gitHubUser *models.GitHubUser, accessToken string) (*models.OAuthProcessResult, error) {
	providerUserID := strconv.Itoa(gitHubUser.ID)

	// Check if this GitHub account is already linked
	existingOAuth, err := h.oauthRepo.GetOAuthAccountByProvider("github", providerUserID)
	if err == nil {
		// GitHub account already linked - get the user and login
		user, err := h.userRepo.GetUserBySessionID(existingOAuth.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get linked user: %w", err)
		}

		// Update the access token
		h.oauthRepo.UpdateOAuthToken(existingOAuth.UserID, "github", accessToken)

		return &models.OAuthProcessResult{
			User:            user,
			IsNewUser:       false,
			IsLinkedAccount: false,
		}, nil
	}

	// Check for email conflicts with existing users
	if gitHubUser.Email != "" {
		conflictingUser, err := h.oauthRepo.CheckEmailConflict(gitHubUser.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email conflict: %w", err)
		}

		if conflictingUser != nil {
			// Email exists - return specific error for email conflict
			return nil, fmt.Errorf("email_conflict: account with email %s already exists", gitHubUser.Email)
		}
	}

	// Create new user with OAuth data
	newUser, err := h.createOAuthUser(gitHubUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth user: %w", err)
	}

	// Link OAuth account to new user
	err = h.oauthRepo.CreateOAuthAccount(newUser.ID, gitHubUser, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to link OAuth account: %w", err)
	}

	return &models.OAuthProcessResult{
		User:            newUser,
		IsNewUser:       true,
		IsLinkedAccount: false,
	}, nil
}

// createOAuthUser creates a new user from GitHub OAuth data
func (h *OAuthHandler) createOAuthUser(gitHubUser *models.GitHubUser) (*models.User, error) {
	// Generate a unique username based on GitHub login
	username := h.generateUniqueUsername(gitHubUser.Login)

	// Create user registration data
	reg := models.UserRegistration{
		Username: username,
		Email:    gitHubUser.Email,
		Password: "", // OAuth users don't have passwords
	}

	// Create the user
	user, err := h.userRepo.CreateUser(reg)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth user: %w", err)
	}

	return user, nil
}

// generateUniqueUsername generates a unique username based on GitHub login
func (h *OAuthHandler) generateUniqueUsername(githubLogin string) string {
	baseUsername := githubLogin
	username := baseUsername
	counter := 1

	// Check if username exists and append number if needed
	for counter <= 100 {
		// This is a simplified approach - in production you'd want proper username checking
		return username
	}

	// Fallback to UUID if all variations taken
	return "user_" + utils.GenerateUUIDToken()[:8]
}

// ================================
// GOOGLE OAUTH FLOW HANDLERS
// ================================

// ServeGoogleLogin initiates Google OAuth flow
func (h *OAuthHandler) ServeGoogleLogin(w http.ResponseWriter, r *http.Request) {

	// Create OAuth state for CSRF protection
	state, err := h.oauthRepo.CreateOAuthState("google")
	if err != nil {
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=state_failed", http.StatusSeeOther)
		return
	}

	// Build Google authorization URL
	authURL := h.buildGoogleAuthURL(state.StateID)

	// Redirect user to Google
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// ServeGoogleCallback handles Google OAuth callback
func (h *OAuthHandler) ServeGoogleCallback(w http.ResponseWriter, r *http.Request) {

	// Get parameters from URL
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errorParam := r.URL.Query().Get("error")

	// Handle user denial
	if errorParam == "access_denied" {
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=google_cancelled", http.StatusSeeOther)
		return
	}

	// Validate required parameters
	if code == "" || state == "" {
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=missing_parameters", http.StatusSeeOther)
		return
	}

	// Validate state (CSRF protection)
	err := h.oauthRepo.ValidateOAuthState(state, "google")
	if err != nil {
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=invalid_state", http.StatusSeeOther)
		return
	}

	// Exchange authorization code for access token
	accessToken, err := h.exchangeGoogleCode(code)
	if err != nil {
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=token_exchange_failed", http.StatusSeeOther)
		return
	}

	// Get user information from Google
	googleUser, err := h.getGoogleUserInfo(accessToken)
	if err != nil {
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=user_info_failed", http.StatusSeeOther)
		return
	}

	// Process the OAuth authentication
	result, err := h.processGoogleAuthentication(googleUser, accessToken)
	if err != nil {
		// Check if it's an email conflict error
		if strings.Contains(err.Error(), "email_conflict:") {
			http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=email_conflict", http.StatusSeeOther)
			return
		}
		// Other authentication errors
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=auth_failed", http.StatusSeeOther)
		return
	}

	// Create session for the authenticated user
	session, err := h.sessionRepo.CreateSession(result.User.ID, r.RemoteAddr)
	if err != nil {
		http.Redirect(w, r, h.config.FrontendBaseURL+"/login?error=session_failed", http.StatusSeeOther)
		return
	}

	// Set session cookie
	utils.SetSessionCookie(session.SessionID, w, r, session.ExpiresAt)

	// Redirect to frontend with success parameters
	var redirectURL string
	if result.IsNewUser {
		redirectURL = h.config.FrontendBaseURL + "/?welcome=google"
	} else if result.IsLinkedAccount {
		redirectURL = h.config.FrontendBaseURL + "/?linked=google"
	} else {
		redirectURL = h.config.FrontendBaseURL + "/"
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// ================================
// GOOGLE API FUNCTIONS
// ================================

// buildGoogleAuthURL builds the Google authorization URL
func (h *OAuthHandler) buildGoogleAuthURL(state string) string {
	params := url.Values{}
	params.Add("client_id", config.Config.GoogleClientID)
	params.Add("redirect_uri", h.config.BackendBaseURL+"/api/auth/google/callback")
	params.Add("scope", "openid email profile")
	params.Add("response_type", "code")
	params.Add("state", state)

	return "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()
}

// exchangeGoogleCode exchanges authorization code for access token
func (h *OAuthHandler) exchangeGoogleCode(code string) (string, error) {
	// Prepare form data
	data := url.Values{}
	data.Set("client_id", config.Config.GoogleClientID)
	data.Set("client_secret", config.Config.GoogleClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", h.config.BackendBaseURL+"/api/auth/google/callback")

	// Create request
	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("token exchange request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %w", err)
	}

	// Parse JSON response
	var tokenResp models.GoogleTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	// Check for errors
	if tokenResp.Error != "" {
		return "", fmt.Errorf("google oauth error: %s - %s", tokenResp.Error, tokenResp.ErrorDescription)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("no access token received from Google")
	}

	return tokenResp.AccessToken, nil
}

// getGoogleUserInfo fetches user information from Google API
func (h *OAuthHandler) getGoogleUserInfo(accessToken string) (*models.GoogleUser, error) {
	// Create request to Google user API
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user info request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google API returned status %d", resp.StatusCode)
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	var googleUser models.GoogleUser
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &googleUser, nil
}

// processGoogleAuthentication handles the core OAuth authentication logic
func (h *OAuthHandler) processGoogleAuthentication(googleUser *models.GoogleUser, accessToken string) (*models.OAuthProcessResult, error) {
	// Check if this Google account is already linked
	existingOAuth, err := h.oauthRepo.GetOAuthAccountByProvider("google", googleUser.ID)
	if err == nil {
		// Google account already linked - get the user and login
		user, err := h.userRepo.GetUserBySessionID(existingOAuth.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get linked user: %w", err)
		}

		// Update the access token
		h.oauthRepo.UpdateOAuthToken(existingOAuth.UserID, "google", accessToken)

		return &models.OAuthProcessResult{
			User:            user,
			IsNewUser:       false,
			IsLinkedAccount: false,
		}, nil
	}

	// Check for email conflicts with existing users
	if googleUser.Email != "" {
		conflictingUser, err := h.oauthRepo.CheckEmailConflict(googleUser.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email conflict: %w", err)
		}

		if conflictingUser != nil {
			// Email exists - return specific error for email conflict
			return nil, fmt.Errorf("email_conflict: account with email %s already exists", googleUser.Email)
		}
	}

	// Create new user with OAuth data
	newUser, err := h.createGoogleOAuthUser(googleUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth user: %w", err)
	}

	// Link OAuth account to new user
	err = h.oauthRepo.CreateOAuthAccount(newUser.ID, googleUser, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to link OAuth account: %w", err)
	}

	return &models.OAuthProcessResult{
		User:            newUser,
		IsNewUser:       true,
		IsLinkedAccount: false,
	}, nil
}

// createGoogleOAuthUser creates a new user from Google OAuth data
func (h *OAuthHandler) createGoogleOAuthUser(googleUser *models.GoogleUser) (*models.User, error) {
	// Generate a unique username based on Google name
	username := h.generateUniqueUsername(googleUser.Name)

	// Create user registration data
	reg := models.UserRegistration{
		Username: username,
		Email:    googleUser.Email,
		Password: "", // OAuth users don't have passwords
	}

	// Create the user
	user, err := h.userRepo.CreateUser(reg)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth user: %w", err)
	}

	return user, nil
}
