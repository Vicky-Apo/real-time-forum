package handlers

import (
	"encoding/json"
	"net/http"

	"frontend-service/internal/models"
	"frontend-service/internal/services"
	"frontend-service/internal/utils"
)

type NotificationHandler struct {
	authService         *services.AuthService
	notificationService *services.NotificationService
	templateService     *services.TemplateService
	errorHandler        *SimpleErrorHandler
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(authService *services.AuthService, notificationService *services.NotificationService, templateService *services.TemplateService) *NotificationHandler {
	return &NotificationHandler{
		authService:         authService,
		notificationService: notificationService,
		templateService:     templateService,
		errorHandler:        NewSimpleErrorHandler(templateService),
	}
}

// ServeNotifications handles the notifications page request (GET /notifications)
func (h *NotificationHandler) ServeNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler.ShowError(w, "Method Not Allowed", "This method is not allowed.")
		return
	}

	// Check if user is logged in
	user := utils.GetUserFromSession(r, h.authService)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get session cookie for API calls
	sessionCookie, err := utils.GetSessionCookie(r, h.authService)
	if err != nil {
		h.errorHandler.ShowError(w, "Authentication Error", "Please log in again.")
		return
	}

	// Get notifications from backend
	notifications, err := h.notificationService.GetNotifications(sessionCookie)
	if err != nil {
		h.errorHandler.ShowError(w, "Failed to Load Notifications", "We're having trouble loading your notifications right now. Please try again later.")
		return
	}

	// Filter out nil notifications and count unread
	var validNotifications []*models.Notification
	unreadCount := 0

	for _, notification := range notifications {
		if notification != nil {
			validNotifications = append(validNotifications, notification)
			if notification.IsUnread() {
				unreadCount++
			}
		}
	}

	// Prepare data for template
	data := models.NotificationPageData{
		Notifications: validNotifications,
		User:          user,
		UnreadCount:   unreadCount,
		TotalCount:    len(validNotifications),
	}

	// Render the notifications template
	if err := h.templateService.Render(w, "notifications.html", data); err != nil {
		h.errorHandler.ShowError(w, "Failed to Render Notifications Page", "We're having trouble loading the notifications page right now. Please try again later.")
		return
	}
}

// HandleMarkAsRead handles AJAX requests to mark notification as read (POST /notifications/mark-read/{id})
func (h *NotificationHandler) HandleMarkAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondWithJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	// Check if user is logged in
	user := utils.GetUserFromSession(r, h.authService)
	if user == nil {
		h.respondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
		return
	}

	// Get notification ID from URL path
	notificationID := r.PathValue("id")
	if notificationID == "" {
		h.respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Notification ID is required"})
		return
	}

	// Get session cookie for API calls
	sessionCookie, err := utils.GetSessionCookie(r, h.authService)
	if err != nil {
		h.respondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
		return
	}

	// Mark notification as read via backend API
	err = h.notificationService.MarkAsRead(notificationID, sessionCookie)
	if err != nil {
		h.respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to mark notification as read"})
		return
	}

	// Return success response
	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Notification marked as read"})
}

// GetNotificationsAPI handles AJAX requests to get notifications (GET /api/notifications)
func (h *NotificationHandler) GetNotificationsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondWithJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	// Check if user is logged in
	user := utils.GetUserFromSession(r, h.authService)
	if user == nil {
		h.respondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
		return
	}

	// Get session cookie for API calls
	sessionCookie, err := utils.GetSessionCookie(r, h.authService)
	if err != nil {
		h.respondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
		return
	}

	// Get notifications from backend
	notifications, err := h.notificationService.GetNotifications(sessionCookie)
	if err != nil {
		h.respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get notifications"})
		return
	}

	// Filter out nil notifications and count unread
	var validNotifications []*models.Notification
	unreadCount := 0

	for _, notification := range notifications {
		if notification != nil {
			validNotifications = append(validNotifications, notification)
			if notification.IsUnread() {
				unreadCount++
			}
		}
	}

	// Return JSON response for AJAX polling
	response := map[string]interface{}{
		"notifications": validNotifications,
		"unread_count":  unreadCount,
		"total_count":   len(validNotifications),
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// Helper method to respond with JSON
func (h *NotificationHandler) respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
