package handlers

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"Social/pkg/models"
)

// validateRegisterRequest checks if the registration request contains valid data
func validateRegisterRequest(req models.RegisterRequest) error {
	// Validate email format
	if err := validateEmail(req.Email); err != nil {
		return err
	}

	// Validate password complexity
	if err := validatePassword(req.Password); err != nil {
		return err
	}

	// Validate first and last names
	if req.FirstName == "" || len(req.FirstName) < 2 {
		return errors.New("first name must be at least 2 characters long")
	}
	if req.LastName == "" || len(req.LastName) < 2 {
		return errors.New("last name must be at least 2 characters long")
	}

	// Validate date of birth
	if err := validateDateOfBirth(req.DateOfBirth); err != nil {
		return err
	}

	return nil
}

// validateLoginRequest checks if the login request contains valid data
func validateLoginRequest(req models.LoginRequest) error {
	// Validate email format
	if err := validateEmail(req.Email); err != nil {
		return err
	}

	// Ensure password is not empty
	if req.Password == "" {
		return errors.New("password cannot be empty")
	}

	return nil
}

// validateEmail checks if the email format is valid
func validateEmail(email string) error {
	// Simple regex for validating email format
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// validatePassword checks if the password meets the complexity requirements
func validatePassword(password string) error {
	// Ensure password has at least 8 characters
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Check for at least one uppercase letter
	if !containsUppercase(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Check for at least one lowercase letter
	if !containsLowercase(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	// Check for at least one digit
	if !containsDigit(password) {
		return errors.New("password must contain at least one number")
	}

	// Check for at least one special character
	if !containsSpecialChar(password) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// Helper functions for password validation
func containsUppercase(s string) bool {
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			return true
		}
	}
	return false
}

func containsLowercase(s string) bool {
	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			return true
		}
	}
	return false
}

func containsDigit(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
}

func containsSpecialChar(s string) bool {
	specialChars := "@$!%*?&"
	for _, c := range s {
		if strings.ContainsRune(specialChars, c) {
			return true
		}
	}
	return false
}

// validateDateOfBirth checks if the date of birth is a valid date and ensures the user is at least 13 years old
func validateDateOfBirth(dateOfBirth string) error {
	dob, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		return errors.New("invalid date of birth format, expected YYYY-MM-DD")
	}

	// Ensure the user is at least 13 years old
	ageLimit := time.Now().AddDate(-13, 0, 0)
	if dob.After(ageLimit) {
		return errors.New("you must be at least 13 years old to register")
	}

	return nil
}
