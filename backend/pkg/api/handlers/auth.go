package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"Social/pkg/api/middlewares"
	"Social/pkg/db/sqlite"
	"Social/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Insert the user into the database
	_, err = sqlite.DB.Exec(`
		INSERT INTO users (email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		req.Email, hashedPassword, req.FirstName, req.LastName, req.DateOfBirth, req.Avatar, req.Nickname, req.AboutMe)
	if err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Find the user by email
	var user models.User
	err := sqlite.DB.QueryRow("SELECT id, email, password FROM users WHERE email = ?", req.Email).
		Scan(&user.ID, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	// Compare the password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token, err := middlewares.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Send the token to the client
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
