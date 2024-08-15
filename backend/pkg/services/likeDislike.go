package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"fmt"
	"time"
)

// LikePost adds a like to a post
func LikePost(userID, postID int) error {
	like := models.Like{
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now(),
	}

	_, err := db.DB.Exec(`INSERT INTO likes (user_id, post_id, created_at) VALUES (?, ?, ?)`,
		like.UserID, like.PostID, like.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to add like: %w", err)
	}
	return nil
}

// DislikePost adds a dislike to a post
func DislikePost(userID, postID int) error {
	dislike := models.Dislike{
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now(),
	}

	_, err := db.DB.Exec(`INSERT INTO dislikes (user_id, post_id, created_at) VALUES (?, ?, ?)`,
		dislike.UserID, dislike.PostID, dislike.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to add dislike: %w", err)
	}
	return nil
}
