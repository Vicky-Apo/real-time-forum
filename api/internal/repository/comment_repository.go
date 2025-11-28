package repository

import (
	"database/sql"
	"errors"
	"time"

	"real-time-forum/internal/models"
	"real-time-forum/internal/utils"
)

type CommentRepository struct {
	db *sql.DB
}

// NewCommentRepository creates a new CommentRepository
func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (cor *CommentRepository) CreateComment(postID, userID, content string) (*models.CreateCommentResponse, error) {
	return utils.ExecuteInTransactionWithResult(cor.db, func(tx *sql.Tx) (*models.CreateCommentResponse, error) {
		// Check if post exists
		var exists int
		err := tx.QueryRow("SELECT COUNT(*) FROM posts WHERE post_id = ?", postID).Scan(&exists)
		if err != nil {
			return nil, err
		}
		// exists == 0 means post does exist
		// exists == 1 means post does NOT exist
		if exists == 0 {
			return nil, errors.New("post not found")
		}

		// Generate UUID for comment
		commentID := utils.GenerateUUIDToken()
		createdAt := time.Now()

		// Insert comment
		_, err = tx.Exec(
			"INSERT INTO comments (comment_id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)",
			commentID, postID, userID, content, createdAt,
		)
		if err != nil {
			return nil, err
		}

		// Return lightweight response - just ID and timestamp
		return &models.CreateCommentResponse{
			CommentID: commentID,
			CreatedAt: createdAt,
		}, nil
	})
}

func (cor *CommentRepository) UpdateComment(commentID, userID, content string) error {
	return utils.ExecuteInTransaction(cor.db, func(tx *sql.Tx) error {
		// Check if comment exists and user owns it
		var ownerID string
		err := tx.QueryRow("SELECT user_id FROM comments WHERE comment_id = ?", commentID).Scan(&ownerID)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("comment not found")
			}
			return err
		}

		// Check ownership
		if ownerID != userID {
			return errors.New("unauthorized: you can only update your own comments")
		}
		now := time.Now()
		// Update comment content and set updated_at
		_, err = tx.Exec(
			"UPDATE comments SET content = ?, updated_at = ? WHERE comment_id = ?",
			content, now, commentID,
		)
		if err != nil {
			return err
		}

		return nil
	})
}

func (cor *CommentRepository) DeleteComment(commentID, userID string) error {
	return utils.ExecuteInTransaction(cor.db, func(tx *sql.Tx) error {
		// Check if comment exists and user owns it
		var ownerID string
		err := tx.QueryRow("SELECT user_id FROM comments WHERE comment_id = ?", commentID).Scan(&ownerID)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("comment not found")
			}
			return err
		}

		// Check ownership
		if ownerID != userID {
			return errors.New("unauthorized: you can only delete your own comments")
		}

		// Delete comment
		_, err = tx.Exec("DELETE FROM comments WHERE comment_id = ?", commentID)
		if err != nil {
			return err
		}

		return nil
	})
}

// ...
// GET comment or comments
// ...

// Get []*comments for PostID
func (cor *CommentRepository) GetCommentsByPostID(postID string, limit, offset int, userID string, options utils.SortOptions) ([]*models.Comment, error) {

	// Build dynamic query with sorting using unified system - UNCHANGED
	orderClause := utils.BuildOrderClause(options.SortBy, utils.ContentTypeComments)

	query := `
		SELECT
			c.comment_id,
			c.post_id,
			c.user_id,
			u.nickname,
			c.content,
			c.created_at,
			c.updated_at,
			NULL as user_reaction
		FROM comments c
		JOIN users u ON c.user_id = u.user_id
		WHERE c.post_id = ?
		` + orderClause + `
		LIMIT ? OFFSET ?`

	rows, err := cor.db.Query(query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment, err := cor.scanCommentRow(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	// Check for iteration errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// ..
// COUNT comments methods
// ..
func (cor *CommentRepository) GetCommentCountByPost(postID string) (int, error) {
	var count int
	err := cor.db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", postID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ..
// /  Helper method
// ..
// Helper method to scan comment rows (updated to handle both *sql.Rows and *sql.Row)
func (cor *CommentRepository) scanCommentRow(scanner interface{}) (*models.Comment, error) {
	var comment models.Comment
	var userReaction sql.NullInt64
	var updatedAt sql.NullTime

	var err error
	switch s := scanner.(type) {
	case *sql.Row:
		err = s.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Nickname,
			&comment.Content,
			&comment.CreatedAt,
			&updatedAt,
			&userReaction,
		)
	case *sql.Rows:
		err = s.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Nickname,
			&comment.Content,
			&comment.CreatedAt,
			&updatedAt,
			&userReaction,
		)
	default:
		return nil, errors.New("invalid scanner type")
	}

	if err != nil {
		return nil, err
	}

	// Handle UpdatedAt
	if updatedAt.Valid {
		comment.UpdatedAt = &updatedAt.Time
	} else {
		comment.UpdatedAt = nil
	}

	return &comment, nil
}
