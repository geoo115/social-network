package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"database/sql"
	"fmt"
	"time"
)

func CreateFollowRequest(request models.FollowRequest) error {
	request.CreatedAt = time.Now()

	_, err := db.DB.Exec(`INSERT INTO follow_requests (sender_id, recipient_id, status, created_at)
		VALUES (?, ?, ?, ?)`, request.SenderID, request.RecipientID, request.Status, request.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create follow request: %w", err)
	}
	return nil
}

func GetFollowRequest(id int) (models.FollowRequest, error) {
	row := db.DB.QueryRow(`SELECT id, sender_id, recipient_id, status, created_at
		FROM follow_requests WHERE id = ?`, id)

	var request models.FollowRequest
	if err := row.Scan(&request.ID, &request.SenderID, &request.RecipientID, &request.Status, &request.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return request, fmt.Errorf("follow request not found")
		}
		return request, fmt.Errorf("failed to get follow request: %w", err)
	}
	return request, nil
}

func UpdateFollowRequest(id int, status string) error {
	_, err := db.DB.Exec(`UPDATE follow_requests SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("failed to update follow request: %w", err)
	}
	return nil
}

func DeleteFollowRequest(id int) error {
	_, err := db.DB.Exec(`DELETE FROM follow_requests WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete follow request: %w", err)
	}
	return nil
}

func AddFollower(followerID int, followedID int) error {
	_, err := db.DB.Exec("INSERT INTO followers (follower_id, followed_id) VALUES (?, ?)", followerID, followedID)
	if err != nil {
		return fmt.Errorf("failed to add follower: %w", err)
	}
	return nil
}
