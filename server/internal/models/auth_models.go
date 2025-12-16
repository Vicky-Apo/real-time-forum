package models

import (
	"fmt"
	"time"
)

// GITHUB API RESPONSE MODELS 


// GitHubTokenResponse - Response from GitHub token exchange
type GitHubTokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// GitHubUser - User information from GitHub API
type GitHubUser struct {
	ID        int    `json:"id"`         // GitHub user ID (number)
	Login     string `json:"login"`      // GitHub username
	Email     string `json:"email"`      // Primary email (may be null)
	Name      string `json:"name"`       // Display name
	AvatarURL string `json:"avatar_url"` // Profile picture URL
}

// GitHubEmail - Email from GitHub emails API
type GitHubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}


// OAuthUserAccount - Represents oauth_user_accounts table
type OAuthUserAccount struct {
	UserID           string    `json:"user_id"`
	Provider         string    `json:"provider"`          // 'github'
	ProviderUserID   string    `json:"provider_user_id"`  // Provider user ID as string
	ProviderEmail    string    `json:"provider_email"`    // Email from provider
	ProviderUsername string    `json:"provider_username"` // Username from provider
	AccessToken      string    `json:"-"`                 // Never expose in JSON
	CreatedAt        time.Time `json:"created_at"`
}

// OAuthFlowState - Represents oauth_flow_states table
type OAuthFlowState struct {
	StateID   string    `json:"state_id"`
	Provider  string    `json:"provider"` // 'github'
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// OAUTH FLOW RESULT

// OAuthProcessResult - Result of processing OAuth callback
type OAuthProcessResult struct {
	User            *User       `json:"user"`
	IsNewUser       bool        `json:"is_new_user"`
	IsLinkedAccount bool        `json:"is_linked_account"`
	GitHubData      *GitHubUser `json:"-"` // Don't expose in API responses
	AccessToken     string      `json:"-"` // Don't expose in API responses
}

// API RESPONSE MODELS

// OAuthLoginResponse - Response after successful OAuth login
type OAuthLoginResponse struct {
	User         User   `json:"user"`          // The logged-in user
	SessionID    string `json:"session_id"`    // Created session ID
	IsNewUser    bool   `json:"is_new_user"`   // Was this a new account creation?
	LinkedGitHub bool   `json:"linked_github"` // Was GitHub account linked to existing user?
}

// HELPER FUNCTIONS

// ConvertGitHubUserToGeneric - Convert GitHub user to generic format for database
func ConvertGitHubUserToGeneric(user *GitHubUser) (string, string, string) {
	providerUserID := fmt.Sprintf("%d", user.ID)
	providerEmail := user.Email
	providerUsername := user.Login

	return providerUserID, providerEmail, providerUsername
}

// GetProviderAuthURL - Get the authorization URL for GitHub
func GetProviderAuthURL(clientID, redirectURI, state string) string {
	return fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email&state=%s",
		clientID, redirectURI, state)
}

// GetProviderTokenURL - Get the token exchange URL for GitHub
func GetProviderTokenURL() string {
	return "https://github.com/login/oauth/access_token"
}

// GetProviderUserInfoURL - Get the user info URL for GitHub
func GetProviderUserInfoURL() string {
	return "https://api.github.com/user"
}

// GITHUB API RESPONSE MODELS 

// GOOGLE API RESPONSE MODELS 

// GoogleTokenResponse - Response from Google token exchange
type GoogleTokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// GoogleUser - User information from Google API
type GoogleUser struct {
	ID            string `json:"id"`             // Google user ID (string)
	Email         string `json:"email"`          // Primary email
	Name          string `json:"name"`           // Display name
	Picture       string `json:"picture"`        // Profile picture URL
	VerifiedEmail bool   `json:"verified_email"` // Is email verified
	GivenName     string `json:"given_name"`     // First name
	FamilyName    string `json:"family_name"`    // Last name
}

// ConvertGoogleUserToGeneric - Convert Google user to generic format for database
func ConvertGoogleUserToGeneric(user *GoogleUser) (string, string, string) {
	providerUserID := user.ID
	providerEmail := user.Email
	providerUsername := user.Name // Google doesn't have username, use name

	return providerUserID, providerEmail, providerUsername
}
