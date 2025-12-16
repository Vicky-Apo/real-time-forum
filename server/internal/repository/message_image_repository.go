package repository

import (
	"database/sql"
	"errors"

	"real-time-forum/internal/models"
	"real-time-forum/internal/utils"
)

type MessageImageRepository struct {
	db *sql.DB
}

func NewMessageImageRepository(db *sql.DB) *MessageImageRepository {
	return &MessageImageRepository{db: db}
}

func (mir *MessageImageRepository) SaveImageRecord(tx *sql.Tx, messageID, imageID, imageURL, originalFilename string) error {
	query := `
		INSERT INTO message_images (image_id, message_id, image_url, original_filename)
		VALUES (?, ?, ?, ?)
	`
	_, err := tx.Exec(query, imageID, messageID, imageURL, originalFilename)
	return err
}

func (mir *MessageImageRepository) GetImagesForMessage(messageID string) ([]models.MessageImage, error) {
	query := `
        SELECT image_id, message_id, image_url, original_filename, uploaded_at
        FROM message_images
        WHERE message_id = ?
        ORDER BY uploaded_at ASC
    `
	rows, err := mir.db.Query(query, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.MessageImage
	for rows.Next() {
		var img models.MessageImage
		err := rows.Scan(&img.ImageID, &img.MessageID, &img.ImageURL, &img.OriginalFilename, &img.UploadedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return images, nil
}

func (mir *MessageImageRepository) DeleteImageByID(imageID string) (*models.MessageImage, error) {
	return utils.ExecuteInTransactionWithResult(mir.db, func(tx *sql.Tx) (*models.MessageImage, error) {
		// Get image info first (to return info & so you can delete the file from disk)
		var img models.MessageImage
		query := `
			SELECT image_id, message_id, image_url, original_filename, uploaded_at
			FROM message_images
			WHERE image_id = ?
		`
		err := tx.QueryRow(query, imageID).Scan(
			&img.ImageID, &img.MessageID, &img.ImageURL, &img.OriginalFilename, &img.UploadedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("image not found")
			}
			return nil, err
		}

		// Delete the image record
		delQuery := `DELETE FROM message_images WHERE image_id = ?`
		_, err = tx.Exec(delQuery, imageID)
		if err != nil {
			return nil, err
		}

		return &img, nil
	})
}

func (mir *MessageImageRepository) DeleteAllImagesForMessage(messageID string) ([]models.MessageImage, error) {
	return utils.ExecuteInTransactionWithResult(mir.db, func(tx *sql.Tx) ([]models.MessageImage, error) {
		// Get all images for the message
		images, err := mir.GetImagesForMessage(messageID)
		if err != nil {
			return nil, err
		}
		if len(images) == 0 {
			return images, nil // Nothing to delete
		}

		// Delete all images for the message
		delQuery := `DELETE FROM message_images WHERE message_id = ?`
		_, err = tx.Exec(delQuery, messageID)
		if err != nil {
			return nil, err
		}

		return images, nil
	})
}
