package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"database/sql"
	"fmt"
	"time"
)

// CreateGroup creates a new group
func CreateGroup(group models.Group) error {
	now := time.Now()

	_, err := db.DB.Exec(`INSERT INTO groups (creator_id, title, description, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)`, group.CreatorID, group.Title, group.Description, now, now)
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}
	return nil
}

// GetGroup retrieves a specific group by its ID
func GetGroup(groupID int) (models.Group, error) {
	var group models.Group
	row := db.DB.QueryRow(`SELECT id, creator_id, title, description, created_at, updated_at FROM groups WHERE id = ?`, groupID)
	err := row.Scan(&group.ID, &group.CreatorID, &group.Title, &group.Description, &group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return group, fmt.Errorf("group not found")
		}
		return group, fmt.Errorf("failed to get group: %w", err)
	}
	return group, nil
}

// JoinGroup adds a user to a group
func JoinGroup(userID, groupID int) error {
	now := time.Now()

	_, err := db.DB.Exec(`INSERT INTO group_memberships (user_id, group_id, joined_at) 
		VALUES (?, ?, ?)`, userID, groupID, now)
	if err != nil {
		return fmt.Errorf("failed to join group: %w", err)
	}
	return nil
}

// LeaveGroup removes a user from a group
func LeaveGroup(userID, groupID int) error {
	now := time.Now()

	_, err := db.DB.Exec(`UPDATE group_memberships SET left_at = ? WHERE user_id = ? AND group_id = ?`, now, userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to leave group: %w", err)
	}
	return nil
}

// CreateGroupEvent creates a new event in a group
func CreateGroupEvent(event models.GroupEvent) error {
	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	_, err := db.DB.Exec(`INSERT INTO group_events (group_id, title, description, day_time, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		event.GroupID, event.Title, event.Description, event.DayTime, event.CreatedAt, event.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create group event: %w", err)
	}
	return nil
}

// GetGroupEvent retrieves a specific group event by its ID
func GetGroupEvent(eventID int) (models.GroupEvent, error) {
	var event models.GroupEvent
	row := db.DB.QueryRow(`SELECT id, group_id, title, description, day_time, created_at, updated_at FROM group_events WHERE id = ?`, eventID)
	err := row.Scan(&event.ID, &event.GroupID, &event.Title, &event.Description, &event.DayTime, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return event, fmt.Errorf("group event not found")
		}
		return event, fmt.Errorf("failed to get group event: %w", err)
	}
	return event, nil
}

// UpdateGroupEvent updates a group event
func UpdateGroupEvent(eventID int, updatedEvent models.GroupEvent) error {
	_, err := db.DB.Exec(`UPDATE group_events SET title = ?, description = ?, day_time = ?, updated_at = ? WHERE id = ?`,
		updatedEvent.Title, updatedEvent.Description, updatedEvent.DayTime, time.Now(), eventID)
	if err != nil {
		return fmt.Errorf("failed to update group event: %w", err)
	}
	return nil
}

// DeleteGroupEvent removes a group event
func DeleteGroupEvent(eventID int) error {
	_, err := db.DB.Exec(`DELETE FROM group_events WHERE id = ?`, eventID)
	if err != nil {
		return fmt.Errorf("failed to delete group event: %w", err)
	}
	return nil
}
