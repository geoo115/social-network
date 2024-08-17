package services

import (
	"Social/pkg/db"
	"database/sql"
	"fmt"
	"time"
)

// LikePost adds a like to a post
func LikePost(userID, postID int) error {
	// Check if the user already liked the post
	row := db.DB.QueryRow(`SELECT 1 FROM likes WHERE user_id = ? AND post_id = ?`, userID, postID)
	var exists int
	err := row.Scan(&exists)
	if err == nil {
		return fmt.Errorf("post already liked by this user")
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check like status: %w", err)
	}

	// Add the like if it does not exist
	_, err = db.DB.Exec(`INSERT INTO likes (user_id, post_id, created_at) VALUES (?, ?, ?)`,
		userID, postID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add like: %w", err)
	}
	return nil
}

// DislikePost adds a dislike to a post
func DislikePost(userID, postID int) error {
	// Check if the user already disliked the post
	row := db.DB.QueryRow(`SELECT 1 FROM dislikes WHERE user_id = ? AND post_id = ?`, userID, postID)
	var exists int
	err := row.Scan(&exists)
	if err == nil {
		return fmt.Errorf("post already disliked by this user")
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check dislike status: %w", err)
	}

	// Add the dislike if it does not exist
	_, err = db.DB.Exec(`INSERT INTO dislikes (user_id, post_id, created_at) VALUES (?, ?, ?)`,
		userID, postID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add dislike: %w", err)
	}
	return nil
}
