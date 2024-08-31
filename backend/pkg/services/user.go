// services/user.go

package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"database/sql"
	"fmt"
	"time"
)

// FindOrCreateUserByEmail finds a user by email or creates a new one
func FindOrCreateUserByEmail(email, provider string) (models.User, error) {
	var user models.User

	// Attempt to find user in the database
	err := db.DB.QueryRow("SELECT id, email FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email)
	if err == nil {
		return user, nil
	} else if err != sql.ErrNoRows {
		return user, fmt.Errorf("error checking for existing user: %w", err)
	}

	// Create a new user if not found (no password for OAuth users)
	_, err = db.DB.Exec("INSERT INTO users (email, provider, created_at, updated_at) VALUES (?, ?, ?, ?)",
	email, provider, time.Now(), time.Now())
	if err != nil {
		return user, fmt.Errorf("failed to create user: %w", err)
	}

	// Retrieve the newly created user
	err = db.DB.QueryRow("SELECT id, email FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email)
	if err != nil {
		return user, fmt.Errorf("error retrieving newly created user: %w", err)
	}

	return user, nil
}
