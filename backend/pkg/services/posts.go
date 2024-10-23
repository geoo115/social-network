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
	// Inserting into the posts table
	query := `INSERT INTO posts (user_id, content, image, privacy, created_at, updated_at) 
              VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))`

	_, err := db.DB.Exec(query, post.UserID, post.Content, post.Image, post.Privacy)
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

// GetAllPosts fetches all posts from the database
func GetAllPosts() ([]models.Post, error) {
	// Define the SQL query
	query := `
	SELECT id, user_id, content, image, privacy, created_at, updated_at 
	FROM posts
	ORDER BY created_at DESC;
	`

	// Query the database
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	// Initialize an empty slice to store the posts
	var posts []models.Post

	// Iterate over the rows and scan each row into a Post struct
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		// Append each post to the slice
		posts = append(posts, post)
	}

	// Check for errors from iterating over the rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over posts: %w", err)
	}

	// Return the slice of posts
	return posts, nil
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
