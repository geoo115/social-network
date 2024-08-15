package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"fmt"
	"time"
)

// SendMessage sends a message to a user or group
func SendMessage(message models.Chat) error {
	query := `
		INSERT INTO chats (sender_id, recipient_id, group_id, message, is_group, created_at) 
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := db.DB.Exec(query, message.SenderID, message.RecipientID, message.GroupID, message.Message, message.IsGroup, time.Now())
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// GetMessages retrieves messages between two users or in a group
func GetMessages(userID, recipientID int, groupID int) ([]models.Chat, error) {
	var messages []models.Chat

	query := `
		SELECT id, sender_id, recipient_id, group_id, message, is_group, created_at
		FROM chats
		WHERE (sender_id = ? AND recipient_id = ?)
		   OR (sender_id = ? AND recipient_id = ?)
		   OR (group_id = ? AND is_group = ?)
		ORDER BY created_at`

	rows, err := db.DB.Query(query, userID, recipientID, recipientID, userID, groupID, true)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve messages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var msg models.Chat
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.RecipientID, &msg.GroupID, &msg.Message, &msg.IsGroup, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating rows: %w", err)
	}

	return messages, nil
}
