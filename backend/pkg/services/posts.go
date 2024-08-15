package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"database/sql"
	"fmt"
	"time"
)

// CreatePost inserts a new post into the database
func CreatePost(post models.Post) error {
	_, err := db.DB.Exec(`INSERT INTO posts (user_id, content, image, privacy, created_at, updated_at) 
        VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))`,
		post.UserID, post.Content, post.Image, post.Privacy)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

// GetPost retrieves a post by ID
func GetPost(postID int) (models.Post, error) {
	row := db.DB.QueryRow(`SELECT id, user_id, content, image, privacy, created_at, updated_at 
		FROM posts WHERE id = ?`, postID)

	var post models.Post
	err := row.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt, &post.UpdatedAt)
	if err == sql.ErrNoRows {
		return post, fmt.Errorf("post not found")
	} else if err != nil {
		return post, fmt.Errorf("failed to retrieve post: %w", err)
	}

	return post, nil
}

// UpdatePost updates the content of a post
func UpdatePost(postID int, updatedPost models.Post) error {
	_, err := db.DB.Exec(`UPDATE posts SET content = ?, image = ?, privacy = ?, updated_at = ? 
		WHERE id = ?`, updatedPost.Content, updatedPost.Image, updatedPost.Privacy, time.Now(), postID)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

// DeletePost removes a post from the database
func DeletePost(postID int) error {
	_, err := db.DB.Exec(`DELETE FROM posts WHERE id = ?`, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}
