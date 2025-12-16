package middleware

import (
	"context"
	"errors"
	"net/http"

	"real-time-forum/config"
	"real-time-forum/internal/models"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/utils"
)

type AuthMiddleware struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

func NewMiddleware(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

type contextKey string

// Define constants using this type
const (
	userContextKey contextKey = "user"
)

// Authenticate middleware verifies authentication and sets user in context
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session cookie using config session name
		cookie, err := r.Cookie(config.Config.SessionName)
		if err != nil {
			// No cookie means not authenticated - just continue
			next.ServeHTTP(w, r)
			return
		}

		// Validate the session
		session, err := m.sessionRepo.GetBySessionID(cookie.Value)
		if err != nil {
			// Session invalid or expired - clear the cookie
			utils.ClearSessionCookie(w)
			next.ServeHTTP(w, r)
			return
		}

		// Get the user
		user, err := m.userRepo.GetUserBySessionID(session.UserID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Set user in context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth middleware ensures the user is authenticated
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(userContextKey)
		if user == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, errors.New("unauthorized access").Error())
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetCurrentUser returns the authenticated user from the context
func GetCurrentUser(r *http.Request) *models.User {
	userValue := r.Context().Value(userContextKey)

	if userValue == nil {
		return nil
	}

	user, ok := userValue.(*models.User)
	if !ok {
		return nil
	}

	return user
}
