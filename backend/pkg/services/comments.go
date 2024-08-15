package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"database/sql"
	"fmt"
	"time"
)

// CreateComment adds a new comment to a post
func CreateComment(comment models.Comment) error {
	query := `
		INSERT INTO comments (user_id, post_id, content, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)`

	_, err := db.DB.Exec(query, comment.UserID, comment.PostID, comment.Content, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

// GetComment retrieves a specific comment
func GetComment(commentID int) (models.Comment, error) {
	var comment models.Comment

	query := `
		SELECT id, user_id, post_id, content, created_at, updated_at
		FROM comments
		WHERE id = ?`

	row := db.DB.QueryRow(query, commentID)
	err := row.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
	if err == sql.ErrNoRows {
		return comment, fmt.Errorf("comment not found")
	} else if err != nil {
		return comment, fmt.Errorf("failed to retrieve comment: %w", err)
	}

	return comment, nil
}

// UpdateComment updates a comment
func UpdateComment(commentID int, updatedComment models.Comment) error {
	query := `
		UPDATE comments
		SET content = ?, updated_at = ?
		WHERE id = ?`

	_, err := db.DB.Exec(query, updatedComment.Content, time.Now(), commentID)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

// DeleteComment removes a comment from a post
func DeleteComment(commentID int) error {
	query := `
		DELETE FROM comments
		WHERE id = ?`

	_, err := db.DB.Exec(query, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}
