package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"database/sql"
	"fmt"
	"time"
)

func CreateGroup(group models.Group) (int, error) {
	now := time.Now()
	group.CreatedAt = now
	group.UpdatedAt = now

	res, err := db.DB.Exec(`INSERT INTO groups (creator_id, title, description, created_at, updated_at) 
                            VALUES (?, ?, ?, ?, ?)`, group.CreatorID, group.Title, group.Description, group.CreatedAt, group.UpdatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to create group: %w", err)
	}

	groupID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve group ID: %w", err)
	}

	return int(groupID), nil
}

func GetGroup(groupID int) (models.Group, error) {
	var group models.Group
	row := db.DB.QueryRow(`SELECT id, creator_id, title, description, created_at, updated_at 
                           FROM groups WHERE id = ?`, groupID)
	err := row.Scan(&group.ID, &group.CreatorID, &group.Title, &group.Description, &group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return group, fmt.Errorf("group not found")
		}
		return group, fmt.Errorf("failed to get group: %w", err)
	}
	return group, nil
}

func ListGroups(offset, limit int, searchTerm string) ([]models.Group, error) {
	rows, err := db.DB.Query(`SELECT id, creator_id, title, description, created_at, updated_at 
                              FROM groups WHERE title LIKE ? LIMIT ? OFFSET ?`, "%"+searchTerm+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.CreatorID, &group.Title, &group.Description, &group.CreatedAt, &group.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func InviteToGroup(invitation models.GroupInvitation) error {
	now := time.Now()
	invitation.InvitedAt = now

	_, err := db.DB.Exec(`INSERT INTO group_invitations (group_id, inviter_id, invitee_id, status, invited_at) 
                          VALUES (?, ?, ?, ?, ?)`, invitation.GroupID, invitation.InviterID, invitation.InviteeID, invitation.Status, invitation.InvitedAt)
	if err != nil {
		return fmt.Errorf("failed to invite user to group: %w", err)
	}
	return nil
}

func RespondToInvitation(invitationID int, status string) error {
	now := time.Now()

	_, err := db.DB.Exec(`UPDATE group_invitations SET status = ?, responded_at = ? 
                          WHERE id = ?`, status, now, invitationID)
	if err != nil {
		return fmt.Errorf("failed to respond to invitation: %w", err)
	}
	return nil
}

func CreateGroupRequest(request models.GroupRequest) error {
	now := time.Now()
	request.RequestedAt = now

	_, err := db.DB.Exec(`INSERT INTO group_requests (group_id, requester_id, status, requested_at) 
                          VALUES (?, ?, ?, ?)`, request.GroupID, request.RequesterID, request.Status, request.RequestedAt)
	if err != nil {
		return fmt.Errorf("failed to create group request: %w", err)
	}
	return nil
}

func RespondToGroupRequest(requestID int, status string) error {
	now := time.Now()

	_, err := db.DB.Exec(`UPDATE group_requests SET status = ?, responded_at = ? 
                          WHERE id = ?`, status, now, requestID)
	if err != nil {
		return fmt.Errorf("failed to respond to group request: %w", err)
	}
	return nil
}

func CreateGroupEvent(event models.GroupEvent) error {
	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	_, err := db.DB.Exec(`INSERT INTO group_events (group_id, title, description, day_time, created_at, updated_at) 
                          VALUES (?, ?, ?, ?, ?, ?)`,
		event.GroupID, event.Title, event.Description, event.DayTime, event.CreatedAt, event.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create group event: %w", err)
	}
	return nil
}

func RSVPEvent(rsvp models.EventRespond) error {
	now := time.Now()
	rsvp.RespondedAt = now

	_, err := db.DB.Exec(`INSERT INTO event_rsvps (event_id, user_id, status, responded_at) 
                          VALUES (?, ?, ?, ?)`, rsvp.EventID, rsvp.UserID, rsvp.Status, rsvp.RespondedAt)
	if err != nil {
		return fmt.Errorf("failed to RSVP for event: %w", err)
	}
	return nil
}

func GetGroupEvent(groupID, eventID int) (models.GroupEvent, error) {
    var event models.GroupEvent
    row := db.DB.QueryRow(`SELECT id, group_id, title, description, day_time, created_at, updated_at 
                           FROM group_events WHERE group_id = ? AND id = ?`, groupID, eventID)
    err := row.Scan(&event.ID, &event.GroupID, &event.Title, &event.Description, &event.DayTime, &event.CreatedAt, &event.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return event, fmt.Errorf("event not found")
        }
        return event, fmt.Errorf("failed to get group event: %w", err)
    }
    return event, nil
}

func JoinGroup(groupID, userID int) error {
    now := time.Now()

    // Check if the user is already a member of the group
    var exists bool
    err := db.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM group_memberships WHERE group_id = ? AND user_id = ?)`, groupID, userID).Scan(&exists)
    if err != nil {
        return fmt.Errorf("failed to check group membership: %w", err)
    }

    if exists {
        return fmt.Errorf("user is already a member of the group")
    }

    // Insert a new membership record
    _, err = db.DB.Exec(`INSERT INTO group_memberships (group_id, user_id, joined_at) 
                         VALUES (?, ?, ?)`, groupID, userID, now)
    if err != nil {
        return fmt.Errorf("failed to join group: %w", err)
    }

    return nil
}

func LeaveGroup(groupID, userID int) error {
    now := time.Now()

    // Update the membership record to indicate the user left the group
    res, err := db.DB.Exec(`UPDATE group_memberships SET left_at = ? 
                            WHERE group_id = ? AND user_id = ? AND left_at IS NULL`, now, groupID, userID)
    if err != nil {
        return fmt.Errorf("failed to leave group: %w", err)
    }

    // Check if any rows were updated
    affectedRows, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to check affected rows: %w", err)
    }
    if affectedRows == 0 {
        return fmt.Errorf("user is not currently a member of the group")
    }

    return nil
}

