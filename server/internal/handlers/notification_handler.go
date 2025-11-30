package handlers

import (
	"net/http"

	"real-time-forum/internal/middleware"
	"real-time-forum/internal/models"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/utils"
)

// GetNotificationsHandler returns all notifications for the authenticated user
func GetNotificationsHandler(nr *repository.NotificationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get authenticated user
		user := middleware.GetCurrentUser(r)

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