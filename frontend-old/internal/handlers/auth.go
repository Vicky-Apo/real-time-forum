package handlers

import (
	"frontend-service/config"
	"frontend-service/internal/models"
	"frontend-service/internal/services"
	"frontend-service/internal/utils"
	"frontend-service/internal/validations"
	"net/http"
	"time"
)

// Add this to your existing AuthHandler struct
type AuthHandler struct {
	authService     *services.AuthService
	templateService *services.TemplateService
	errorHandler    *SimpleErrorHandler // ADD THIS
	config          *config.Config
}

// Update your NewAuthHandler constructor
func NewAuthHandler(authService *services.AuthService, templateService *services.TemplateService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService:     authService,
		templateService: templateService,
		errorHandler:    NewSimpleErrorHandler(templateService), // ADD THIS
		config:          cfg,
	}
}

// ServeRegister handles GET and POST requests for registration
func (h *AuthHandler) ServeRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.showRegisterForm(w, r)
	case http.MethodPost:
		h.handleRegisterForm(w, r)
	default:
		h.errorHandler.ShowError(w, "Method Not Allowed", "This method is not allowed for registration.")
	}
}

// ServeLogin handles GET and POST requests for login
func (h *AuthHandler) ServeLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.showLoginForm(w, r)
	case http.MethodPost:
		h.handleLoginForm(w, r)
	default:
		h.errorHandler.ShowError(w, "Method Not Allowed", "This method is not allowed for login.")
	}
}

// ServeLogout handles logout requests
func (h *AuthHandler) ServeLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		h.errorHandler.ShowError(w, "Method Not Allowed", "This method is not allowed for logout.")
		return
	}

	// Get session cookie using config value
	cookie, err := utils.GetSessionCookie(r, h.authService)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Call backend to logout
	if err := h.authService.LogoutUser(cookie.Value); err != nil {
	} else {
	}

	// Clear the session cookie on frontend using config value
	utils.ClearSessionCookie(h.config.SessionName, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// showRegisterForm displays the registration form (GET request)
func (h *AuthHandler) showRegisterForm(w http.ResponseWriter, _ *http.Request) {
	data := models.RegisterPageData{
		FormData: &models.UserRegistration{}, // Empty form data for initial load
	}

	if err := h.templateService.Render(w, "register.html", data); err != nil {
		h.errorHandler.ShowError(w, "Failed to Render Register Page", "We're having trouble loading the registration page right now. Please try again later.")
	}
}

// handleRegisterForm processes the registration form submission (POST request)
func (h *AuthHandler) handleRegisterForm(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		h.showRegisterError(w, "Invalid form data", &models.UserRegistration{})
		return
	}

	// Get form values - ADD confirm_password here
	formData := models.UserRegistration{
		Username:        r.FormValue("username"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm_password"),
	}

	// Basic validation (check if all fields are provided)
	if formData.Username == "" || formData.Email == "" || formData.Password == "" {
		h.showRegisterError(w, "All fields are required", &formData)
		return
	}

	// NEW: Check password confirmation
	if formData.ConfirmPassword == "" {
		h.showRegisterError(w, "Password confirmation is required", &formData)
		return
	}

	// NEW: Validate passwords match (frontend validation)
	if formData.Password != formData.ConfirmPassword {
		h.showRegisterError(w, "Passwords do not match", &formData)
		return
	}

	// FRONTEND VALIDATION using your validation functions
	if err := validations.ValidateUserInput(formData.Username, formData.Email, formData.Password); err != nil {
		h.showRegisterError(w, err.Error(), &formData)
		return
	}

	// Call backend API to register user (backend will also validate)
	if err := h.authService.RegisterUser(formData); err != nil {
		h.showRegisterError(w, err.Error(), &formData)
		return
	}

	// Registration successful - show success message with empty form
	data := models.RegisterPageData{
		Success:  "Registration successful! You can now login.",
		FormData: &models.UserRegistration{}, // Clear form on success
	}

	if err := h.templateService.Render(w, "login.html", data); err != nil {
		h.errorHandler.ShowError(w, "Failed to Render Register Page", "We're having trouble loading the registration page right now. Please try again later.")
	}
}

// showRegisterError displays registration form with error message AND preserved form data
func (h *AuthHandler) showRegisterError(w http.ResponseWriter, errorMsg string, formData *models.UserRegistration) {
	// Clear password for security - user will need to retype it
	if formData != nil {
		formData.Password = ""
	}

	data := models.RegisterPageData{
		Error:    errorMsg,
		FormData: formData, // Pass back the form data so fields stay populated
	}

	if err := h.templateService.Render(w, "register.html", data); err != nil {
		h.errorHandler.ShowError(w, "Failed to Render Register Page", "We're having trouble loading the registration page right now. Please try again later.")
	}
}

// showLoginForm displays the login form (GET request)
func (h *AuthHandler) showLoginForm(w http.ResponseWriter, r *http.Request) {
	// Check for OAuth error parameter in URL
	errorParam := r.URL.Query().Get("error")
	if errorParam != "" {
		// Handle OAuth errors with nice error pages
		switch errorParam {
		case "email_conflict":
			h.errorHandler.ShowOAuthError(w,
				"An account with this email address already exists. Please log in with your existing email and password. After logging in, you can connect your GitHub account from your profile settings.")
		case "github_cancelled":
			h.errorHandler.ShowOAuthError(w,
				"GitHub authorization was cancelled. You can try again or log in with your email and password.")
		case "auth_failed":
			h.errorHandler.ShowOAuthError(w,
				"GitHub authentication failed. Please try again or log in with your email and password.")
		case "token_exchange_failed":
			h.errorHandler.ShowOAuthError(w,
				"Failed to connect with GitHub. Please try again later or use email/password login.")
		case "user_info_failed":
			h.errorHandler.ShowOAuthError(w,
				"Could not get your information from GitHub. Please try again or use email/password login.")
		case "session_failed":
			h.errorHandler.ShowOAuthError(w,
				"Could not create your session. Please try again.")
		case "state_failed":
			h.errorHandler.ShowOAuthError(w,
				"Security validation failed. Please try the GitHub login again.")
		default:
			h.errorHandler.ShowOAuthError(w,
				"Login failed. Please try again or use email and password.")
		}
		return
	}

	// Check for success parameters (welcome messages)
	welcomeParam := r.URL.Query().Get("welcome")
	if welcomeParam == "github" {
		// User successfully registered with GitHub - redirect to home with success
		http.Redirect(w, r, "/?success=github_welcome", http.StatusSeeOther)
		return
	}

	linkedParam := r.URL.Query().Get("linked")
	if linkedParam == "github" {
		// User successfully linked GitHub - redirect to home with success
		http.Redirect(w, r, "/?success=github_linked", http.StatusSeeOther)
		return
	}

	// Normal login form rendering
	data := models.LoginPageData{
		FormData: &models.UserLogin{}, // Empty form data for initial load
	}

	if err := h.templateService.Render(w, "login.html", data); err != nil {
		h.errorHandler.ShowError(w, "Failed to Render Login Page", "We're having trouble loading the login page right now. Please try again later.")
	}
}

// handleLoginForm processes the login form submission (POST request)
func (h *AuthHandler) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		h.showLoginError(w, "Invalid form data", &models.UserLogin{})
		return
	}

	// Get form values
	formData := models.UserLogin{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	// Basic validation (check if all fields are provided)
	if formData.Email == "" || formData.Password == "" {
		h.showLoginError(w, "Email and password are required", &formData)
		return
	}

	// FRONTEND VALIDATION using your validation functions
	err := validations.ValidateEmail(formData.Email)
	if err != nil {
		h.showLoginError(w, err.Error(), &formData)
		return
	}
	err = validations.ValidatePassword(formData.Password)
	if err != nil {
		h.showLoginError(w, err.Error(), &formData)
		return
	}

	// Call backend API to login user
	_, sessionID, err := h.authService.LoginUser(formData)
	if err != nil {

		h.showLoginError(w, err.Error(), &formData)
		return
	}

	// Set session cookie using config value and utility function
	expiresAt := time.Now().Add(24 * time.Hour)
	utils.SetSessionCookie(h.config.SessionName, sessionID, w, r, expiresAt) // CHANGED: Use utility with config session name

	// Login successful - redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// showLoginError displays login form with error message AND preserved form data
func (h *AuthHandler) showLoginError(w http.ResponseWriter, errorMsg string, formData *models.UserLogin) {
	// Clear password for security - user will need to retype it
	if formData != nil {
		formData.Password = ""
	}

	data := models.LoginPageData{
		Error:    errorMsg,
		FormData: formData, // Pass back the form data so fields stay populated
	}

	if err := h.templateService.Render(w, "login.html", data); err != nil {
		h.errorHandler.ShowError(w, "fail to render template", "try again later")
	}
}
