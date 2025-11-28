package repository

import (
	"database/sql"
	"errors"
	"time"

	"platform.zone01.gr/git/gpapadopoulos/forum/internal/models"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/utils"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) CreateUser(reg models.UserRegistration) (*models.User, error) {
	return utils.ExecuteInTransactionWithResult(ur.DB, func(tx *sql.Tx) (*models.User, error) {
		// Check if username exists
		var usernameCount int
		err := tx.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", reg.Username).Scan(&usernameCount)
		if err != nil {
			return nil, err
		}
		if usernameCount > 0 {
			return nil, errors.New("username already taken")
		}

		// Check if email exists
		var emailCount int
		err = tx.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", reg.Email).Scan(&emailCount)
		if err != nil {
			return nil, err
		}
		if emailCount > 0 {
			return nil, errors.New("email already taken")
		}

		userID := utils.GenerateUUIDToken()
		createdAt := time.Now()

		// Hash the password
		hashedPassword, err := utils.HashPassword(reg.Password)
		if err != nil {
			return nil, err
		}

		// Insert user record
		_, err = tx.Exec(
			"INSERT INTO users (user_id, username, email, password_hash, created_at) VALUES (?, ?, ?, ?, ?)",
			userID, reg.Username, reg.Email, hashedPassword, createdAt,
		)
		if err != nil {
			return nil, err
		}

		// Return the created user
		return &models.User{
			ID:        userID,
			Username:  reg.Username,
			Email:     reg.Email,
			CreatedAt: createdAt,
		}, nil
	})
}

// GetBySessionID retrieves a user by session ID
func (ur *UserRepository) GetUserBySessionID(id string) (*models.User, error) {
	var user models.User

	err := ur.DB.QueryRow(
		"SELECT user_id, username, email, created_at FROM users WHERE user_id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil

}

// GetAuthByUserID retrieves user authentication data by user ID
func (ur *UserRepository) GetAuthByUserID(userID string) (*models.UserPassword, error) {
	var auth models.UserPassword

	err := ur.DB.QueryRow(
		"SELECT user_id, password_hash FROM users WHERE user_id = ?",
		userID,
	).Scan(&auth.UserID, &auth.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user authentication not found")
		}
		return nil, err
	}

	return &auth, nil

}

// GetUserByEmail retrieves a user by email
func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := ur.DB.QueryRow(
		"SELECT user_id, username, email, created_at FROM users WHERE LOWER(email) = LOWER(?)",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil

}

// Authenticate validates a user's login credentials
func (ur *UserRepository) Authenticate(login models.UserLogin) (*models.User, error) {
	// Get the user by email
	user, err := ur.GetUserByEmail(login.Email)
	if err != nil {
		return nil, errors.New("email not found")
	}

	// Get the user's authentication data
	auth, err := ur.GetAuthByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	// Check the password
	if !utils.CheckPasswordHash(login.Password, auth.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (ur *UserRepository) GetCurrentUser(userID string) (*models.User, error) {

	var user models.User

	err := ur.DB.QueryRow(
		"SELECT user_id, username, email, created_at FROM users WHERE user_id = ?",
		userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil

}

// GetUserProfile retrieves complete user profile with statistics
func (ur *UserRepository) GetUserProfile(userID string) (*models.UserProfile, error) {
	// First get basic user info using existing method
	user, err := ur.GetCurrentUser(userID)
	if err != nil {
		return nil, err
	}

	// Get profile statistics
	stats, err := ur.GetProfileStats(userID)
	if err != nil {
		return nil, err
	}

	profile := &models.UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Stats:     *stats,
	}

	return profile, nil
}

// GetProfileStats calculates all profile statistics for a user
func (ur *UserRepository) GetProfileStats(userID string) (*models.ProfileStats, error) {
	stats := &models.ProfileStats{}

	// 1. Count total posts by user
	err := ur.DB.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&stats.TotalPosts)
	if err != nil {
		return nil, err
	}

	// 2. Count total comments by user
	err = ur.DB.QueryRow("SELECT COUNT(*) FROM comments WHERE user_id = ?", userID).Scan(&stats.TotalComments)
	if err != nil {
		return nil, err
	}

	// 3. Count posts this user has liked - FIXED: use post_reactions table
	err = ur.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM post_reactions 
		WHERE user_id = ? AND reaction_type = 1
	`, userID).Scan(&stats.PostsLiked)
	if err != nil {
		return nil, err
	}

	// 4. Count posts this user has commented on (NEW)
	err = ur.DB.QueryRow(`
		SELECT COUNT(DISTINCT p.post_id) 
		FROM posts p
		JOIN comments c ON p.post_id = c.post_id
		WHERE c.user_id = ?
	`, userID).Scan(&stats.PostsCommentedOn)
	if err != nil {
		return nil, err
	}

	// 5. Count total likes received (on posts + comments) - FIXED: use separate tables
	err = ur.DB.QueryRow(`
		SELECT (
			SELECT COALESCE(COUNT(*), 0) FROM post_reactions pr 
			JOIN posts p ON pr.post_id = p.post_id 
			WHERE pr.reaction_type = 1 AND p.user_id = ?
		) + (
			SELECT COALESCE(COUNT(*), 0) FROM comment_reactions cr 
			JOIN comments c ON cr.comment_id = c.comment_id 
			WHERE cr.reaction_type = 1 AND c.user_id = ?
		)
	`, userID, userID).Scan(&stats.LikesReceived)
	if err != nil {
		return nil, err
	}

	// 6. Count total dislikes received (on posts + comments) - FIXED: use separate tables
	err = ur.DB.QueryRow(`
		SELECT (
			SELECT COALESCE(COUNT(*), 0) FROM post_reactions pr 
			JOIN posts p ON pr.post_id = p.post_id 
			WHERE pr.reaction_type = 2 AND p.user_id = ?
		) + (
			SELECT COALESCE(COUNT(*), 0) FROM comment_reactions cr 
			JOIN comments c ON cr.comment_id = c.comment_id 
			WHERE cr.reaction_type = 2 AND c.user_id = ?
		)
	`, userID, userID).Scan(&stats.DislikesReceived)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
