package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"real-time-forum/internal/models"
	"real-time-forum/internal/utils"
)

type OAuthRepository struct {
	DB *sql.DB
}

// NewOAuthRepository creates a new OAuth repository
func NewOAuthRepository(db *sql.DB) *OAuthRepository {
	return &OAuthRepository{DB: db}
}

// ================================
// OAUTH FLOW STATES OPERATIONS
// ================================

// CreateOAuthState creates a new OAuth flow state for CSRF protection
func (r *OAuthRepository) CreateOAuthState(provider string) (*models.OAuthFlowState, error) {
	// Generate secure random state ID
	stateID := utils.GenerateUUIDToken()
	now := time.Now()
	expiresAt := now.Add(15 * time.Minute) // States expire in 15 minutes

	// Insert state into database
	_, err := r.DB.Exec(`
		INSERT INTO oauth_flow_states (state_id, provider, created_at, expires_at)
		VALUES (?, ?, ?, ?)`,
		stateID, provider, now, expiresAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth state: %w", err)
	}

	return &models.OAuthFlowState{
		StateID:   stateID,
		Provider:  provider,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}, nil
}

// ValidateOAuthState validates and consumes an OAuth state
func (r *OAuthRepository) ValidateOAuthState(stateID, provider string) error {
	return utils.ExecuteInTransaction(r.DB, func(tx *sql.Tx) error {
		// Check if state exists and is valid
		var dbProvider string
		var expiresAt time.Time

		err := tx.QueryRow(`
			SELECT provider, expires_at 
			FROM oauth_flow_states 
			WHERE state_id = ?`,
			stateID).Scan(&dbProvider, &expiresAt)

		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("invalid OAuth state")
			}
			return fmt.Errorf("failed to validate OAuth state: %w", err)
		}

		// Check if state expired
		if time.Now().After(expiresAt) {
			return errors.New("OAuth state expired")
		}

		// Check if provider matches
		if dbProvider != provider {
			return errors.New("OAuth state provider mismatch")
		}

		// Delete the state (one-time use)
		_, err = tx.Exec("DELETE FROM oauth_flow_states WHERE state_id = ?", stateID)
		if err != nil {
			return fmt.Errorf("failed to consume OAuth state: %w", err)
		}

		return nil
	})
}

// ================================
// OAUTH USER ACCOUNTS OPERATIONS
// ================================

// GetOAuthAccountByProvider gets OAuth account by provider and provider user ID
func (r *OAuthRepository) GetOAuthAccountByProvider(provider, providerUserID string) (*models.OAuthUserAccount, error) {
	var account models.OAuthUserAccount

	err := r.DB.QueryRow(`
		SELECT user_id, provider, provider_user_id, provider_email, provider_username, access_token, created_at
		FROM oauth_user_accounts 
		WHERE provider = ? AND provider_user_id = ?`,
		provider, providerUserID).Scan(
		&account.UserID, &account.Provider, &account.ProviderUserID,
		&account.ProviderEmail, &account.ProviderUsername,
		&account.AccessToken, &account.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("OAuth account not found")
		}
		return nil, fmt.Errorf("failed to get OAuth account: %w", err)
	}

	return &account, nil
}

// CreateOAuthAccount creates a new OAuth account link (UPDATED for both GitHub and Google)
func (r *OAuthRepository) CreateOAuthAccount(userID string, userData interface{}, accessToken string) error {
	return utils.ExecuteInTransaction(r.DB, func(tx *sql.Tx) error {
		var provider, providerUserID, providerEmail, providerUsername string

		// Handle different provider types
		switch user := userData.(type) {
		case *models.GitHubUser:
			provider = "github"
			providerUserID, providerEmail, providerUsername = models.ConvertGitHubUserToGeneric(user)
		case *models.GoogleUser:
			provider = "google"
			providerUserID, providerEmail, providerUsername = models.ConvertGoogleUserToGeneric(user)
		default:
			return fmt.Errorf("unsupported OAuth provider type")
		}

		// Check if this OAuth account is already linked to another user
		var existingUserID string
		err := tx.QueryRow(`
			SELECT user_id FROM oauth_user_accounts 
			WHERE provider = ? AND provider_user_id = ?`,
			provider, providerUserID).Scan(&existingUserID)

		if err == nil {
			// OAuth account already linked
			if existingUserID != userID {
				return fmt.Errorf("%s account already linked to different user", provider)
			}
			// Already linked to same user - update token
			_, err = tx.Exec(`
				UPDATE oauth_user_accounts 
				SET access_token = ?, provider_email = ?, provider_username = ? 
				WHERE user_id = ? AND provider = ?`,
				accessToken, providerEmail, providerUsername, userID, provider)
			return err
		} else if err != sql.ErrNoRows {
			return fmt.Errorf("failed to check existing OAuth account: %w", err)
		}

		// Create new OAuth account link
		_, err = tx.Exec(`
			INSERT INTO oauth_user_accounts 
			(user_id, provider, provider_user_id, provider_email, provider_username, access_token, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			userID, provider, providerUserID, providerEmail, providerUsername, accessToken, time.Now())

		if err != nil {
			return fmt.Errorf("failed to create OAuth account: %w", err)
		}

		return nil
	})
}

// UpdateOAuthToken updates the access token for an OAuth account
func (r *OAuthRepository) UpdateOAuthToken(userID, provider, newToken string) error {
	result, err := r.DB.Exec(`
		UPDATE oauth_user_accounts 
		SET access_token = ? 
		WHERE user_id = ? AND provider = ?`,
		newToken, userID, provider)

	if err != nil {
		return fmt.Errorf("failed to update OAuth token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("OAuth account not found")
	}

	return nil
}

// ================================
// EMAIL CONFLICT CHECKING
// ================================

// CheckEmailConflict checks if an email from OAuth conflicts with existing users
func (r *OAuthRepository) CheckEmailConflict(email string) (*models.User, error) {
	if email == "" {
		return nil, nil // No email, no conflict
	}

	var user models.User
	err := r.DB.QueryRow(`
		SELECT user_id, username, email, created_at
		FROM users 
		WHERE LOWER(email) = LOWER(?)`,
		email).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No conflict
		}
		return nil, fmt.Errorf("failed to check email conflict: %w", err)
	}

	return &user, nil // Conflict found
}
