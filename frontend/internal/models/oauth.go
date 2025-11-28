package models

// // ================================
// // OAUTH CALLBACK MODELS (Essential for frontend)
// // ================================

// // OAuthCallbackData - Data received from backend OAuth callback
// type OAuthCallbackData struct {
// 	User         *User  `json:"user"`          // The authenticated user
// 	SessionID    string `json:"session_id"`    // Session ID from backend
// 	IsNewUser    bool   `json:"is_new_user"`   // Was this a new account creation?
// 	LinkedGitHub bool   `json:"linked_github"` // Was GitHub account linked to existing user?
// 	Error        string `json:"error"`         // Error if OAuth failed
// }

// // OAuthResult - Result of processing OAuth callback on frontend
// type OAuthResult struct {
// 	Success     bool   `json:"success"`
// 	User        *User  `json:"user,omitempty"`
// 	SessionID   string `json:"session_id,omitempty"`
// 	IsNewUser   bool   `json:"is_new_user"`
// 	Error       string `json:"error,omitempty"`
// 	RedirectURL string `json:"redirect_url,omitempty"` // Where to redirect after OAuth
// }
