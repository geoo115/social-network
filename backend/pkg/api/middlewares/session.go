package middlewares

import (
	"Social/pkg/db"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"
)

// GenerateSessionID generates a new session ID
func GenerateSessionID(userID int) (string, error) {
	sessionID := generateRandomString(32)
	expiresAt := time.Now().Add(24 * time.Hour) // Sessions expire after 24 hours

	_, err := db.DB.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", sessionID, userID, expiresAt)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// RetrieveSession retrieves session data by session ID
func RetrieveSession(sessionID string) (int, error) {
	var userID int
	var expiresAt time.Time

	err := db.DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_id = ?", sessionID).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err
	}

	if time.Now().After(expiresAt) {
		return 0, sql.ErrNoRows // Session expired
	}

	return userID, nil
}

// DeleteSession deletes a session by session ID
func DeleteSession(sessionID string) error {
	_, err := db.DB.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}

// generateRandomString generates a random string of the given length
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
