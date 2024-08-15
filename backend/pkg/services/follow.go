package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"database/sql"
	"fmt"
	"time"
)

// CreateFollowRequest creates a new follow request in the database
func CreateFollowRequest(request models.FollowRequest) error {
	request.CreatedAt = time.Now()

	_, err := db.DB.Exec(`INSERT INTO follow_requests (sender_id, recipient_id, status, created_at)
		VALUES (?, ?, ?, ?)`, request.SenderID, request.RecipientID, request.Status, request.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create follow request: %w", err)
	}
	return nil
}

// GetFollowRequest retrieves a follow request by ID
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

// UpdateFollowRequest updates the status of a follow request
func UpdateFollowRequest(id int, status string) error {
	_, err := db.DB.Exec(`UPDATE follow_requests SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("failed to update follow request: %w", err)
	}
	return nil
}

// DeleteFollowRequest deletes a follow request by ID
func DeleteFollowRequest(id int) error {
	_, err := db.DB.Exec(`DELETE FROM follow_requests WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete follow request: %w", err)
	}
	return nil
}
