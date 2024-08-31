package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// getCurrentUserID retrieves the user ID from the request context.
func getCurrentUserID(r *http.Request) (int, error) {
	userIDInterface := r.Context().Value("userID")
	if userIDInterface == nil {
		return 0, errors.New("user ID not found in context")
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		return 0, errors.New("invalid user ID type in context")
	}

	return userID, nil
}

// GetProfile handles retrieving a user's profile information.
func GetProfile(w http.ResponseWriter, r *http.Request, userIDStr string) {
	// Parse the requested user ID
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Retrieve the current logged-in user ID from the context
	requesterID, err := getCurrentUserID(r)
	if err != nil {
		http.Error(w, "Unable to retrieve current user ID", http.StatusInternalServerError)
		return
	}

	// Retrieve profile information
	profile, posts, followers, following, err := services.GetProfile(requesterID, userID)
	if err != nil {
		switch err.Error() {
		case "profile not found":
			http.Error(w, "Profile not found", http.StatusNotFound)
		case "profile is private and you are not a follower":
			http.Error(w, "Profile is private and you are not a follower", http.StatusForbidden)
		default:
			http.Error(w, "Failed to retrieve profile", http.StatusInternalServerError)
		}
		return
	}

	// Create response structure
	response := struct {
		User      models.User   `json:"user"`
		Posts     []models.Post `json:"posts"`
		Followers []models.User `json:"followers"`
		Following []models.User `json:"following"`
	}{
		User:      profile,
		Posts:     posts,
		Followers: followers,
		Following: following,
	}

	// Encode response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode profile", http.StatusInternalServerError)
	}
}

// UpdateProfile handles updating a user's profile information.
func UpdateProfile(w http.ResponseWriter, r *http.Request, userIDStr string) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the user ID from the request URL
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Update the profile
	if err := services.UpdateProfile(userID, user); err != nil {
		http.Error(w, "Failed to update profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}
