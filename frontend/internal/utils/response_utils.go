// internal/utils/response_utils.go
package utils

import (
	"net/http"
)

// RenderTemplate - Render template with error handling
func RenderTemplate(templateService TemplateService, w http.ResponseWriter, templateName string, data interface{}) error {
	if err := templateService.Render(w, templateName, data); err != nil {
		return NewGeneralError("Failed to render page", http.StatusInternalServerError)
	}
	return nil
}

// TemplateService interface
type TemplateService interface {
	Render(w http.ResponseWriter, templateName string, data interface{}) error
}

// RedirectToPost - Redirect to specific post
func RedirectToPost(w http.ResponseWriter, r *http.Request, postID string) {
	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}

// RedirectToHome - Redirect to home page
func RedirectToHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// RedirectToLogin - Redirect to login page
func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// SetSuccessMessage - Set success message in session/cookie
func SetSuccessMessage(w http.ResponseWriter, message string) {
	cookie := &http.Cookie{
		Name:     "success_message",
		Value:    message,
		Path:     "/",
		MaxAge:   300,   // 5 minutes
		Secure:   false, // Set to true in production with HTTPS
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}
