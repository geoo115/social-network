package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"Social/pkg/models"
	"Social/pkg/services"
)

// GetProfile retrieves a user's profile by ID
func GetProfile(w http.ResponseWriter, r *http.Request, userIDStr string) {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	profile, err := services.GetProfile(userID)
	if err != nil {
		if err.Error() == "profile not found" {
			http.Error(w, "Profile not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve profile", http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(profile); err != nil {
		http.Error(w, "Failed to encode profile", http.StatusInternalServerError)
	}
}

// UpdateProfile updates a user's profile information
func UpdateProfile(w http.ResponseWriter, r *http.Request, userIDStr string) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = services.UpdateProfile(userID, user)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
