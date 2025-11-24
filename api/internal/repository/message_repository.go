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

// GetConversations returns all users organized by last message timestamp
// Users with messages are sorted by most recent first
// Users without messages are sorted alphabetically at the end
func (mr *MessageRepository) GetConversations(currentUserID string) ([]models.Conversation, error) {
	query := `
		WITH user_conversations AS (
			SELECT
				CASE
					WHEN m.sender_id = ? THEN m.recipient_id
					ELSE m.sender_id
				END AS other_user_id,
				MAX(m.created_at) AS last_message_time
			FROM messages m
			WHERE m.sender_id = ? OR m.recipient_id = ?
			GROUP BY other_user_id
		),
		last_message_details AS (
			SELECT
				uc.other_user_id,
				uc.last_message_time,
				m.content AS last_message_content,
				m.created_at AS last_message_created_at,
				CASE WHEN m.sender_id = ? THEN 1 ELSE 0 END AS is_from_me
			FROM user_conversations uc
			JOIN messages m ON (
				((m.sender_id = ? AND m.recipient_id = uc.other_user_id)
				OR (m.sender_id = uc.other_user_id AND m.recipient_id = ?))
				AND m.created_at = uc.last_message_time
			)
		),
		unread_counts AS (
			SELECT
				sender_id AS other_user_id,
				COUNT(*) AS unread_count
			FROM messages
			WHERE recipient_id = ? AND is_read = 0
			GROUP BY sender_id
		)
		SELECT
			u.user_id,
			u.nickname,
			lmd.last_message_content,
			lmd.last_message_created_at,
			lmd.is_from_me,
			COALESCE(uc.unread_count, 0) AS unread_count
		FROM users u
		LEFT JOIN last_message_details lmd ON u.user_id = lmd.other_user_id
		LEFT JOIN unread_counts uc ON u.user_id = uc.other_user_id
		WHERE u.user_id != ?
		ORDER BY
			CASE WHEN lmd.last_message_time IS NULL THEN 1 ELSE 0 END,
			lmd.last_message_time DESC,
			u.nickname ASC
	`

	rows, err := mr.db.Query(query,
		currentUserID, currentUserID, currentUserID,
		currentUserID, currentUserID, currentUserID,
		currentUserID, currentUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conversations := []models.Conversation{}
	for rows.Next() {
		var conv models.Conversation
		var lastMsgContent sql.NullString
		var lastMsgCreatedAt sql.NullTime
		var isFromMe sql.NullBool

		err := rows.Scan(
			&conv.UserID,
			&conv.Nickname,
			&lastMsgContent,
			&lastMsgCreatedAt,
			&isFromMe,
			&conv.UnreadCount,
		)
		if err != nil {
			return nil, err
		}

		// Only populate LastMessage if there is one
		if lastMsgContent.Valid && lastMsgCreatedAt.Valid && isFromMe.Valid {
			conv.LastMessage = &models.LastMessage{
				Content:   lastMsgContent.String,
				CreatedAt: lastMsgCreatedAt.Time,
				IsFromMe:  isFromMe.Bool,
			}
		} else {
			conv.LastMessage = nil
		}

		conversations = append(conversations, conv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}
