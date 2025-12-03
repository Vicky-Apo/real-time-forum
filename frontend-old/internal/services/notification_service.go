// internal/services/notification_service.go
package services

import (
	"net/http"

	"frontend-service/internal/models"
	"frontend-service/internal/utils"
)

type NotificationService struct {
	*BaseClient
}

// NewNotificationService creates a new notification service
func NewNotificationService(baseClient *BaseClient) *NotificationService {
	return &NotificationService{
		BaseClient: baseClient,
	}
}

// GetNotifications fetches all notifications for the authenticated user
func (s *NotificationService) GetNotifications(sessionCookie *http.Cookie) ([]*models.Notification, error) {
	// Make API request using utils
	apiResponse, err := utils.MakeGETRequest(s.HTTPClient, s.BaseURL, "/notifications", sessionCookie)
	if err != nil {
		return nil, err
	}

	// Convert response to notification response
	var notificationResponse models.NotificationResponse
	if err := utils.ConvertAPIData(apiResponse.Data, &notificationResponse); err != nil {
		return nil, utils.NewGeneralError("Failed to parse notification response", 500)
	}

	return notificationResponse.Notifications, nil
}

// MarkAsRead marks a specific notification as read
func (s *NotificationService) MarkAsRead(notificationID string, sessionCookie *http.Cookie) error {
	// Make API request using utils
	_, err := utils.MakePOSTRequest(s.HTTPClient, s.BaseURL, "/notifications/mark-read/"+notificationID, nil, sessionCookie)
	return err
}

// GetUnreadCount gets the count of unread notifications for the user
func (s *NotificationService) GetUnreadCount(sessionCookie *http.Cookie) (int, error) {
	// Make API request using utils
	apiResponse, err := utils.MakeGETRequest(s.HTTPClient, s.BaseURL, "/notifications/unread-count", sessionCookie)
	if err != nil {
		return 0, err
	}

	// Convert response to count
	var countResponse struct {
		UnreadCount int `json:"unread_count"`
	}

	if err := utils.ConvertAPIData(apiResponse.Data, &countResponse); err != nil {
		return 0, utils.NewGeneralError("Failed to parse unread count", 500)
	}

	return countResponse.UnreadCount, nil
}

// MarkAllAsRead marks all notifications as read for the authenticated user
func (s *NotificationService) MarkAllAsRead(sessionCookie *http.Cookie) error {
	// Make API request using utils
	_, err := utils.MakePOSTRequest(s.HTTPClient, s.BaseURL, "/notifications/mark-all-read", nil, sessionCookie)
	return err
}
