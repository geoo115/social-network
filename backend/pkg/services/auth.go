package services

import (
	"Social/pkg/db"
	"Social/pkg/models"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// RegisterUser creates a new user in the database
func RegisterUser(user models.RegisterRequest) error {
	var existingUser models.User
	err := db.DB.QueryRow("SELECT id FROM users WHERE email = ?", user.Email).Scan(&existingUser.ID)
	if err == nil {
		return fmt.Errorf("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	dateOfBirth, err := time.Parse("2006-01-02", user.DateOfBirth)
	if err != nil {
		return fmt.Errorf("invalid date of birth format: %w", err)
	}

	_, err = db.DB.Exec(`INSERT INTO users (email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Email,
		hashedPassword,
		user.FirstName,
		user.LastName,
		dateOfBirth,
		user.Avatar,
		user.Nickname,
		user.AboutMe,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}

// AuthenticateUser verifies user credentials
func AuthenticateUser(email, password string) (models.User, error) {
	var user models.User
	row := db.DB.QueryRow("SELECT id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, created_at, updated_at FROM users WHERE email = ?", email)
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
			return user, fmt.Errorf("user not found: %w", err)
		}
		return user, fmt.Errorf("error retrieving user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, fmt.Errorf("invalid password: %w", err)
	}

	return user, nil
}
