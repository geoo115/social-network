package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// getCurrentUserID extracts the current user ID from the request context.
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

func GetProfile(w http.ResponseWriter, r *http.Request, userIDStr string) {
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

	profile, posts, followers, following, err := services.GetProfile(requesterID, userID)
	if err != nil {
		if err.Error() == "profile not found" {
			http.Error(w, "Profile not found", http.StatusNotFound)
		} else if err.Error() == "profile is private and you are not a follower" {
			http.Error(w, "Profile is private and you are not a follower", http.StatusForbidden)
		} else {
			http.Error(w, "Failed to retrieve profile", http.StatusInternalServerError)
		}
		return
	}

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

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode profile", http.StatusInternalServerError)
	}
}

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
