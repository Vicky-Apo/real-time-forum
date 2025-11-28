package handlers

import (
	"net/http"

	"platform.zone01.gr/git/gpapadopoulos/forum/internal/middleware"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/models"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/repository"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/utils"
)

// GetNotificationsHandler returns all notifications for the authenticated user
func GetNotificationsHandler(nr *repository.NotificationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)
		if user == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// Get all notifications for the user
		notifications, err := nr.GetUserNotifications(user.ID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve notifications")
			return
		}

		// Create response
		response := models.NotificationResponse{
			Notifications: notifications,
			TotalCount:    len(notifications),
		}

		utils.RespondWithSuccess(w, http.StatusOK, response)
	}
}

// MarkAsReadHandler marks a specific notification as read
func MarkAsReadHandler(nr *repository.NotificationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)
		if user == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// Get notification ID from URL path
		notificationID := r.PathValue("id")
		if notificationID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Notification ID is required")
			return
		}

		// Mark notification as read
		err := nr.MarkAsRead(notificationID, user.ID)
		if err != nil {
			if err.Error() == "notification not found or not owned by user" {
				utils.RespondWithError(w, http.StatusNotFound, "Notification not found")
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to mark notification as read")
			return
		}

		utils.RespondWithSuccess(w, http.StatusOK, map[string]string{"message": "Notification marked as read"})
	}
}
// RegisterNotificationRoutes registers the notification routes