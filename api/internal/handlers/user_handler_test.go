package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"real-time-forum/config"
	"real-time-forum/internal/models"
	"real-time-forum/internal/repository"
)

// init runs before all tests to initialize config
func init() {
	// Load config with defaults (no .env file needed for tests)
	config.LoadConfig()
}

// MockUserRepository is a test implementation of UserRepositoryInterface
type MockUserRepository struct {
	CreateUserFunc              func(models.UserRegistration) (*models.User, error)
	GetUserBySessionIDFunc      func(string) (*models.User, error)
	GetUserByNicknameOrEmailFunc func(string) (*models.User, error)
	GetAuthByUserIDFunc         func(string) (*models.UserPassword, error)
	AuthenticateFunc            func(models.UserLogin) (*models.User, error)
	GetCurrentUserFunc          func(string) (*models.User, error)
}

// Implement the interface methods
func (m *MockUserRepository) CreateUser(reg models.UserRegistration) (*models.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(reg)
	}
	return nil, errors.New("not implemented")
}

func (m *MockUserRepository) GetUserBySessionID(id string) (*models.User, error) {
	if m.GetUserBySessionIDFunc != nil {
		return m.GetUserBySessionIDFunc(id)
	}
	return nil, errors.New("not implemented")
}

func (m *MockUserRepository) GetUserByNicknameOrEmail(identifier string) (*models.User, error) {
	if m.GetUserByNicknameOrEmailFunc != nil {
		return m.GetUserByNicknameOrEmailFunc(identifier)
	}
	return nil, errors.New("not implemented")
}

func (m *MockUserRepository) GetAuthByUserID(userID string) (*models.UserPassword, error) {
	if m.GetAuthByUserIDFunc != nil {
		return m.GetAuthByUserIDFunc(userID)
	}
	return nil, errors.New("not implemented")
}

func (m *MockUserRepository) Authenticate(login models.UserLogin) (*models.User, error) {
	if m.AuthenticateFunc != nil {
		return m.AuthenticateFunc(login)
	}
	return nil, errors.New("not implemented")
}

func (m *MockUserRepository) GetCurrentUser(userID string) (*models.User, error) {
	if m.GetCurrentUserFunc != nil {
		return m.GetCurrentUserFunc(userID)
	}
	return nil, errors.New("not implemented")
}

// Compile-time check that mock implements the interface
var _ repository.UserRepositoryInterface = (*MockUserRepository)(nil)

// ===== ACTUAL TESTS =====

// TestRegisterHandler_EmailAlreadyTaken tests the "email already taken" error path
// BEFORE interfaces: Required database setup, test data, cleanup - SLOW
// AFTER interfaces: Runs in microseconds, no database needed - FAST
func TestRegisterHandler_EmailAlreadyTaken(t *testing.T) {
	// Create a mock that simulates "email already taken"
	mockRepo := &MockUserRepository{
		CreateUserFunc: func(reg models.UserRegistration) (*models.User, error) {
			return nil, errors.New("email already taken")
		},
	}

	// Create the handler with our mock
	handler := RegisterHandler(mockRepo)

	// Create a test registration request
	registration := models.UserRegistration{
		Nickname:        "testuser",
		Email:           "test@example.com",
		Password:        "ValidPass123!",
		ConfirmPassword: "ValidPass123!",
		FirstName:       "Test",
		LastName:        "User",
		Age:             25,
		Gender:          "Male", // Must be capitalized
	}

	body, _ := json.Marshal(registration)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	// Execute the handler
	handler(rec, req)

	// Assert the response
	if rec.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, rec.Code)
	}

	// Verify the error message
	var response map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&response)

	if response["error"] != "email already taken" {
		t.Errorf("Expected error 'email already taken', got '%v'", response["error"])
	}
}

// TestRegisterHandler_NicknameAlreadyTaken tests the "nickname already taken" error path
func TestRegisterHandler_NicknameAlreadyTaken(t *testing.T) {
	mockRepo := &MockUserRepository{
		CreateUserFunc: func(reg models.UserRegistration) (*models.User, error) {
			return nil, errors.New("nickname already taken")
		},
	}

	handler := RegisterHandler(mockRepo)

	registration := models.UserRegistration{
		Nickname:        "existinguser",
		Email:           "newemail@example.com",
		Password:        "ValidPass123!",
		ConfirmPassword: "ValidPass123!",
		FirstName:       "Test",
		LastName:        "User",
		Age:             25,
		Gender:          "Male",
	}

	body, _ := json.Marshal(registration)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, rec.Code)
	}
}

// TestRegisterHandler_Success tests successful registration
func TestRegisterHandler_Success(t *testing.T) {
	// Mock that returns a successful user creation
	mockRepo := &MockUserRepository{
		CreateUserFunc: func(reg models.UserRegistration) (*models.User, error) {
			return &models.User{
				ID:        "test-user-id-123",
				Nickname:  reg.Nickname,
				Email:     reg.Email,
				FirstName: reg.FirstName,
				LastName:  reg.LastName,
				Age:       reg.Age,
				Gender:    reg.Gender,
			}, nil
		},
	}

	handler := RegisterHandler(mockRepo)

	registration := models.UserRegistration{
		Nickname:        "newuser",
		Email:           "newuser@example.com",
		Password:        "ValidPass123!",
		ConfirmPassword: "ValidPass123!",
		FirstName:       "New",
		LastName:        "User",
		Age:             30,
		Gender:          "Female",
	}

	body, _ := json.Marshal(registration)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	// Verify the returned user
	var response map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&response)

	data := response["data"].(map[string]interface{})
	if data["nickname"] != "newuser" {
		t.Errorf("Expected nickname 'newuser', got '%v'", data["nickname"])
	}
}

// TestRegisterHandler_MissingFields tests validation of required fields
func TestRegisterHandler_MissingFields(t *testing.T) {
	// Mock is never called because validation fails first
	mockRepo := &MockUserRepository{}
	handler := RegisterHandler(mockRepo)

	// Registration with missing email
	registration := models.UserRegistration{
		Nickname:        "testuser",
		Email:           "", // Missing!
		Password:        "ValidPass123!",
		ConfirmPassword: "ValidPass123!",
		FirstName:       "Test",
		LastName:        "User",
		Age:             25,
		Gender:          "Male",
	}

	body, _ := json.Marshal(registration)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
