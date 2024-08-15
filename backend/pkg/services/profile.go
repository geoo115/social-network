package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"database/sql"
	"fmt"
	"time"
)

// GetProfile retrieves a user's profile information
func GetProfile(userID int) (models.User, error) {
	var user models.User
	row := db.DB.QueryRow(`
		SELECT id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, created_at, updated_at
		FROM users WHERE id = ?`, userID)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.Avatar,
		&user.Nickname,
		&user.AboutMe,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("profile not found")
		}
		return user, fmt.Errorf("failed to get profile: %w", err)
	}
	return user, nil
}

// UpdateProfile updates user profile information
func UpdateProfile(userID int, updatedProfile models.User) error {
	_, err := db.DB.Exec(`UPDATE users SET first_name = ?, last_name = ?, date_of_birth = ?, avatar = ?, nickname = ?, about_me = ?, updated_at = ? 
		WHERE id = ?`, updatedProfile.FirstName, updatedProfile.LastName, updatedProfile.DateOfBirth, updatedProfile.Avatar, updatedProfile.Nickname, updatedProfile.AboutMe, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	return nil
}
