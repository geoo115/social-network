package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"database/sql"
	"fmt"
	"time"
)

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

// CreateComment adds a new comment to a post
func CreateComment(comment models.Comment) error {
	// Optional: Check if the comment already exists
	row := db.DB.QueryRow(`SELECT 1 FROM comments WHERE user_id = ? AND post_id = ? AND content = ?`,
		comment.UserID, comment.PostID, comment.Content)
	var exists int
	err := row.Scan(&exists)
	if err == nil {
		return fmt.Errorf("duplicate comment detected")
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check comment existence: %w", err)
	}

	// Add the comment if it does not exist
	_, err = db.DB.Exec(`INSERT INTO comments (user_id, post_id, content, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)`,
		comment.UserID, comment.PostID, comment.Content, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}

// UpdateComment updates a comment
func UpdateComment(commentID int, updatedComment models.Comment) error {
	// Check if the comment exists before updating
	row := db.DB.QueryRow(`SELECT 1 FROM comments WHERE id = ?`, commentID)
	var exists int
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return fmt.Errorf("comment not found")
	}
	if err != nil {
		return fmt.Errorf("failed to check comment existence: %w", err)
	}

	_, err = db.DB.Exec(`UPDATE comments SET content = ?, updated_at = ? WHERE id = ?`,
		updatedComment.Content, time.Now(), commentID)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	return nil
}

// DeleteComment removes a comment from a post
func DeleteComment(commentID int) error {
	// Check if the comment exists before deleting
	row := db.DB.QueryRow(`SELECT 1 FROM comments WHERE id = ?`, commentID)
	var exists int
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return fmt.Errorf("comment not found")
	}
	if err != nil {
		return fmt.Errorf("failed to check comment existence: %w", err)
	}

	_, err = db.DB.Exec(`DELETE FROM comments WHERE id = ?`, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}
