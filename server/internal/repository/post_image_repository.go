package repository

import (
	"database/sql"
	"errors"

	"platform.zone01.gr/git/gpapadopoulos/forum/internal/models"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/utils"
)

type PostImagesRepository struct {
	db *sql.DB
}

func NewPostImagesRepository(db *sql.DB) *PostImagesRepository {
	return &PostImagesRepository{db: db}
}

func (pir *PostImagesRepository) SaveImageRecord(postID, imageID, imageURL, originalFilename string) error {
	return utils.ExecuteInTransaction(pir.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO post_images (image_id, post_id, image_url, original_filename)
			VALUES (?, ?, ?, ?)
		`
		_, err := tx.Exec(query, imageID, postID, imageURL, originalFilename)
		return err
	})
}

func (pir *PostImagesRepository) GetImagesForPost(postID string) ([]models.PostImage, error) {
	query := `
        SELECT image_id, post_id, image_url, original_filename, uploaded_at
        FROM post_images
        WHERE post_id = ?
        ORDER BY uploaded_at ASC
    `
	rows, err := pir.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.PostImage
	for rows.Next() {
		var img models.PostImage
		err := rows.Scan(&img.ImageID, &img.PostID, &img.ImageURL, &img.OriginalFilename, &img.UploadedAt)
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

func (pir *PostImagesRepository) DeleteImageByID(imageID string) (*models.PostImage, error) {
	return utils.ExecuteInTransactionWithResult(pir.db, func(tx *sql.Tx) (*models.PostImage, error) {
		// Get image info first (to return info & so you can delete the file from disk)
		var img models.PostImage
		query := `
			SELECT image_id, post_id, image_url, original_filename, uploaded_at
			FROM post_images
			WHERE image_id = ?
		`
		err := tx.QueryRow(query, imageID).Scan(
			&img.ImageID, &img.PostID, &img.ImageURL, &img.OriginalFilename, &img.UploadedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("image not found")
			}
			return nil, err
		}

		// Delete the image record
		delQuery := `DELETE FROM post_images WHERE image_id = ?`
		_, err = tx.Exec(delQuery, imageID)
		if err != nil {
			return nil, err
		}

		return &img, nil
	})
}

func (pir *PostImagesRepository) DeleteAllImagesForPost(postID string) ([]models.PostImage, error) {
	return utils.ExecuteInTransactionWithResult(pir.db, func(tx *sql.Tx) ([]models.PostImage, error) {
		// Get all images for the post
		images, err := pir.GetImagesForPost(postID)
		if err != nil {
			return nil, err
		}
		if len(images) == 0 {
			return images, nil // Nothing to delete
		}

		// Delete all images for the post
		delQuery := `DELETE FROM post_images WHERE post_id = ?`
		_, err = tx.Exec(delQuery, postID)
		if err != nil {
			return nil, err
		}

		return images, nil
	})
}
