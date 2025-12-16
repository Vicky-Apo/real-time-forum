package utils

import (
	"net/http"
	"time"

	"real-time-forum/config"
)

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     config.Config.SessionName, // Use config session name
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,                // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode, // Consistent with frontend
	})
}

func SetSessionCookie(value string, w http.ResponseWriter, r *http.Request, expiresAt time.Time) {
	// Set the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     config.Config.SessionName, // Use config session name
		Value:    value,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false,                // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode, // Consistent with frontend
	})
}
