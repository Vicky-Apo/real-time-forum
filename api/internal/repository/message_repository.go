package repository

import (
	"database/sql"
	"errors"
	"time"

	"real-time-forum/internal/models"
	"real-time-forum/internal/utils"
)

type MessageRepository struct {
	db *sql.DB
}

// NewMessageRepository creates a new MessageRepository
func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// SaveMessage saves a new message to the database
func (mr *MessageRepository) SaveMessage(senderID, recipientID, content string) (*models.SendMessageResponse, error) {
	return utils.ExecuteInTransactionWithResult(mr.db, func(tx *sql.Tx) (*models.SendMessageResponse, error) {
		// Check if recipient exists
		var exists int
		err := tx.QueryRow("SELECT COUNT(*) FROM users WHERE user_id = ?", recipientID).Scan(&exists)
		if err != nil {
			return nil, err
		}
		if exists == 0 {
			return nil, errors.New("recipient not found")
		}

		// Generate UUID for message
		messageID := utils.GenerateUUIDToken()
		createdAt := time.Now()

		// Insert message
		_, err = tx.Exec(
			"INSERT INTO messages (message_id, sender_id, recipient_id, content, created_at, is_read) VALUES (?, ?, ?, ?, ?, ?)",
			messageID, senderID, recipientID, content, createdAt, false,
		)
		if err != nil {
			return nil, err
		}

		// Return lightweight response
		return &models.SendMessageResponse{
			MessageID: messageID,
			CreatedAt: createdAt,
		}, nil
	})
}

// GetMessages retrieves message history between the current user and another user
// Returns messages in descending order (newest first) for initial load
func (mr *MessageRepository) GetMessages(currentUserID, otherUserID string, limit int, beforeTimestamp *time.Time) (*models.GetMessagesResponse, error) {
	var rows *sql.Rows
	var err error

	// Build query based on whether we have a beforeTimestamp (for pagination)
	baseQuery := `
		SELECT m.message_id, m.sender_id, u.nickname, m.recipient_id, m.content, m.created_at, m.is_read
		FROM messages m
		JOIN users u ON m.sender_id = u.user_id
		WHERE ((m.sender_id = ? AND m.recipient_id = ?) OR (m.sender_id = ? AND m.recipient_id = ?))
	`

	if beforeTimestamp != nil {
		// Pagination: get messages before the specified timestamp
		query := baseQuery + ` AND m.created_at < ? ORDER BY m.created_at DESC LIMIT ?`
		rows, err = mr.db.Query(query, currentUserID, otherUserID, otherUserID, currentUserID, beforeTimestamp, limit+1)
	} else {
		// Initial load: get the most recent messages
		query := baseQuery + ` ORDER BY m.created_at DESC LIMIT ?`
		rows, err = mr.db.Query(query, currentUserID, otherUserID, otherUserID, currentUserID, limit+1)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []models.Message{}
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(
			&msg.MessageID,
			&msg.SenderID,
			&msg.SenderName,
			&msg.RecipientID,
			&msg.Content,
			&msg.CreatedAt,
			&msg.IsRead,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Check if there are more messages
	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit] // Trim to requested limit
	}

	return &models.GetMessagesResponse{
		Messages: messages,
		HasMore:  hasMore,
	}, nil
}

// MarkMessagesAsRead marks all messages from a specific sender to the current user as read
func (mr *MessageRepository) MarkMessagesAsRead(currentUserID, senderID string) error {
	return utils.ExecuteInTransaction(mr.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(
			"UPDATE messages SET is_read = 1 WHERE recipient_id = ? AND sender_id = ? AND is_read = 0",
			currentUserID, senderID,
		)
		return err
	})
}

// GetUnreadCount gets the count of unread messages for a user
func (mr *MessageRepository) GetUnreadCount(userID string) (int, error) {
	var count int
	err := mr.db.QueryRow(
		"SELECT COUNT(*) FROM messages WHERE recipient_id = ? AND is_read = 0",
		userID,
	).Scan(&count)
	return count, err
}
