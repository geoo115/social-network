package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func GetProfile(requesterID, userID int) (models.User, []models.Post, []models.User, []models.User, error) {
	var user models.User

	row := db.DB.QueryRow(`
        SELECT id, email, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_private, created_at, updated_at
        FROM users WHERE id = ?`, userID)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.Avatar,
		&user.Nickname,
		&user.AboutMe,
		&user.IsPrivate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil, nil, nil, fmt.Errorf("profile not found")
		}
		log.Printf("Error retrieving profile: %v", err) // Log the detailed error
		return user, nil, nil, nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// Profile visibility check
	if user.IsPrivate && requesterID != userID {
		row := db.DB.QueryRow(`
            SELECT 1 FROM followers WHERE follower_id = ? AND followed_id = ?`, requesterID, userID)
		err = row.Scan()
		if err == sql.ErrNoRows {
			return user, nil, nil, nil, fmt.Errorf("profile is private and you are not a follower")
		}
		if err != nil {
			log.Printf("Error checking follow status: %v", err) // Log the detailed error
			return user, nil, nil, nil, fmt.Errorf("failed to check follow status: %w", err)
		}
	}

	// Fetch posts, followers, and following
	posts, err := fetchPosts(userID)
	if err != nil {
		log.Printf("Error fetching posts: %v", err) // Log the detailed error
		return user, nil, nil, nil, fmt.Errorf("failed to get posts: %w", err)
	}

	followers, err := fetchFollowers(userID)
	if err != nil {
		log.Printf("Error fetching followers: %v", err) // Log the detailed error
		return user, nil, nil, nil, fmt.Errorf("failed to get followers: %w", err)
	}

	following, err := fetchFollowing(userID)
	if err != nil {
		log.Printf("Error fetching following: %v", err) // Log the detailed error
		return user, nil, nil, nil, fmt.Errorf("failed to get following: %w", err)
	}

	return user, posts, followers, following, nil
}

func fetchPosts(userID int) ([]models.Post, error) {
	rows, err := db.DB.Query(`SELECT id, user_id, content, image, privacy, created_at, updated_at FROM posts WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}

		// Fetch comments for the post
		comments, err := fetchComments(post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments

		// Fetch likes count for the post
		likes, err := fetchLikes(post.ID)
		if err != nil {
			return nil, err
		}
		post.Likes = likes

		// Fetch dislikes count for the post
		dislikes, err := fetchDislikes(post.ID)
		if err != nil {
			return nil, err
		}
		post.Dislikes = dislikes

		posts = append(posts, post)
	}
	return posts, nil
}

func fetchComments(postID int) ([]models.Comment, error) {
	rows, err := db.DB.Query(`SELECT id, user_id, post_id, content, created_at, updated_at FROM comments WHERE post_id = ?`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func fetchLikes(postID int) (int, error) {
	var count int
	err := db.DB.QueryRow(`SELECT COUNT(*) FROM likes WHERE post_id = ?`, postID).Scan(&count)
	return count, err
}

func fetchDislikes(postID int) (int, error) {
	var count int
	err := db.DB.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE post_id = ?`, postID).Scan(&count)
	return count, err
}

func fetchFollowers(userID int) ([]models.User, error) {
	rows, err := db.DB.Query(`SELECT u.id, u.email, u.first_name, u.last_name, u.date_of_birth, u.avatar, u.nickname, u.about_me, u.is_private, u.created_at, u.updated_at
        FROM followers f JOIN users u ON f.follower_id = u.id WHERE f.followed_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.AboutMe, &user.IsPrivate, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		followers = append(followers, user)
	}
	return followers, nil
}

func fetchFollowing(userID int) ([]models.User, error) {
	rows, err := db.DB.Query(`SELECT u.id, u.email, u.first_name, u.last_name, u.date_of_birth, u.avatar, u.nickname, u.about_me, u.is_private, u.created_at, u.updated_at
        FROM followers f JOIN users u ON f.followed_id = u.id WHERE f.follower_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.AboutMe, &user.IsPrivate, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		following = append(following, user)
	}
	return following, nil
}

func UpdateProfile(userID int, updatedProfile models.User) error {
	_, err := db.DB.Exec(`UPDATE users SET first_name = ?, last_name = ?, date_of_birth = ?, avatar = ?, nickname = ?, about_me = ?, is_private = ?, updated_at = ? 
        WHERE id = ?`, updatedProfile.FirstName, updatedProfile.LastName, updatedProfile.DateOfBirth, updatedProfile.Avatar, updatedProfile.Nickname, updatedProfile.AboutMe, updatedProfile.IsPrivate, time.Now(), userID)
	if err != nil {
		log.Printf("Error updating profile: %v", err) // Log the detailed error
		return fmt.Errorf("failed to update profile: %w", err)
	}
	return nil
}
