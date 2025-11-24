package repository

import (
	"database/sql"
	"errors"
	"time"

	"real-time-forum/internal/models"
	"real-time-forum/internal/utils"
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
		// Check if nickname exists
		var nicknameCount int
		err := tx.QueryRow("SELECT COUNT(*) FROM users WHERE nickname = ?", reg.Nickname).Scan(&nicknameCount)
		if err != nil {
			return nil, err
		}
		if nicknameCount > 0 {
			return nil, errors.New("nickname already taken")
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

		// Insert user record with new fields
		_, err = tx.Exec(
			"INSERT INTO users (user_id, nickname, age, gender, first_name, last_name, email, password_hash, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			userID, reg.Nickname, reg.Age, reg.Gender, reg.FirstName, reg.LastName, reg.Email, hashedPassword, createdAt,
		)
		if err != nil {
			return nil, err
		}

		// Return the created user
		return &models.User{
			ID:        userID,
			Nickname:  reg.Nickname,
			Age:       reg.Age,
			Gender:    reg.Gender,
			FirstName: reg.FirstName,
			LastName:  reg.LastName,
			Email:     reg.Email,
			CreatedAt: createdAt,
		}, nil
	})
}

// GetBySessionID retrieves a user by session ID
func (ur *UserRepository) GetUserBySessionID(id string) (*models.User, error) {
	var user models.User

	err := ur.DB.QueryRow(
		"SELECT user_id, nickname, age, gender, first_name, last_name, email, created_at FROM users WHERE user_id = ?",
		id,
	).Scan(&user.ID, &user.Nickname, &user.Age, &user.Gender, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
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

// GetUserByNicknameOrEmail retrieves a user by nickname or email
func (ur *UserRepository) GetUserByNicknameOrEmail(identifier string) (*models.User, error) {
	var user models.User

	err := ur.DB.QueryRow(
		"SELECT user_id, nickname, age, gender, first_name, last_name, email, created_at FROM users WHERE LOWER(nickname) = LOWER(?) OR LOWER(email) = LOWER(?)",
		identifier, identifier,
	).Scan(&user.ID, &user.Nickname, &user.Age, &user.Gender, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil

}

// Authenticate validates a user's login credentials (accepts nickname or email)
func (ur *UserRepository) Authenticate(login models.UserLogin) (*models.User, error) {
	// Get the user by nickname or email
	user, err := ur.GetUserByNicknameOrEmail(login.Identifier)
	if err != nil {
		return nil, errors.New("invalid credentials")
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
		"SELECT user_id, nickname, age, gender, first_name, last_name, email, created_at FROM users WHERE user_id = ?",
		userID,
	).Scan(&user.ID, &user.Nickname, &user.Age, &user.Gender, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil

}
