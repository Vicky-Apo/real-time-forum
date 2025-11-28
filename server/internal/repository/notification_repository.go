package repository

import (
	"database/sql"
	"fmt"

	"platform.zone01.gr/git/gpapadopoulos/forum/internal/models"
)

type NotificationRepository struct {
	DB *sql.DB
}

// NewNotificationRepository creates a new NotificationRepository
func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

// CreateNotification creates a new notification in the database
func (nr *NotificationRepository) CreateNotification(notification *models.Notification) error {
	query := `
		INSERT INTO notifications (notification_id, user_id, trigger_username, post_content_preview, post_id, action, is_read, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := nr.DB.Exec(query,
		notification.NotificationID,
		notification.UserID,
		notification.TriggerUsername,
		notification.PostContentPreview,
		notification.PostID,
		notification.Action,
		notification.IsRead,
		notification.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}

// GetUserNotifications retrieves all notifications for a specific user, ordered by newest first
func (nr *NotificationRepository) GetUserNotifications(userID string) ([]*models.Notification, error) {
	query := `
		SELECT notification_id, user_id, trigger_username, post_content_preview, post_id, action, is_read, created_at
		FROM notifications
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := nr.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		notification := &models.Notification{}
		err := rows.Scan(
			&notification.NotificationID,
			&notification.UserID,
			&notification.TriggerUsername,
			&notification.PostContentPreview,
			&notification.PostID,
			&notification.Action,
			&notification.IsRead,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notifications: %w", err)
	}

	return notifications, nil
}

// MarkAsRead marks a specific notification as read for a user
func (nr *NotificationRepository) MarkAsRead(notificationID, userID string) error {
	query := `
		UPDATE notifications 
		SET is_read = TRUE 
		WHERE notification_id = ? AND user_id = ?
	`

	result, err := nr.DB.Exec(query, notificationID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification not found or not owned by user")
	}

	return nil
}
