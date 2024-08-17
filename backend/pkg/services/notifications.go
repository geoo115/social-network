package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"fmt"
)

func GetNotifications(userID int) ([]models.Notification, error) {
	var notifications []models.Notification

	rows, err := db.DB.Query(`SELECT id, user_id, type, message, is_read, created_at, details FROM notifications WHERE user_id = ?`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var notification models.Notification
		if err := rows.Scan(&notification.ID, &notification.UserID, &notification.Type, &notification.Message, &notification.IsRead, &notification.CreatedAt, &notification.Details); err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over notifications: %w", err)
	}

	return notifications, nil
}

func CreateNotification(notification models.Notification) error {
	if notification.Type == "" || notification.Message == "" {
		return fmt.Errorf("notification type or message cannot be empty")
	}

	query := `
		INSERT INTO notifications (user_id, type, message, is_read, created_at, details)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := db.DB.Exec(query, notification.UserID, notification.Type, notification.Message, notification.IsRead, notification.CreatedAt, notification.Details)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}

func MarkNotificationAsRead(notificationID int) error {
	query := `UPDATE notifications SET is_read = true WHERE id = ?`
	_, err := db.DB.Exec(query, notificationID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	return nil
}
