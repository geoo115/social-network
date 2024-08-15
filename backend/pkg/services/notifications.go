package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"fmt"
)

// GetNotifications retrieves notifications for a user
func GetNotifications(userID int) ([]models.Notification, error) {
	var notifications []models.Notification

	rows, err := db.DB.Query(`SELECT id, user_id, message, is_read, created_at FROM notifications WHERE user_id = ?`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var notification models.Notification
		if err := rows.Scan(&notification.ID, &notification.UserID, &notification.Message, &notification.IsRead, &notification.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over notifications: %w", err)
	}

	return notifications, nil
}
