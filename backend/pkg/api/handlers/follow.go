package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"encoding/json"
	"net/http"
	"strconv"
)

// CreateFollowRequest handles the creation of a new follow request
func CreateFollowRequest(w http.ResponseWriter, r *http.Request) {
	var request models.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	request.SenderID = userID
	request.Status = models.FollowRequestPending

	if err := services.CreateFollowRequest(request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetFollowRequest retrieves a follow request by ID
func GetFollowRequest(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid follow request ID", http.StatusBadRequest)
		return
	}

	followRequest, err := services.GetFollowRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(followRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateFollowRequest updates the status of an existing follow request
func UpdateFollowRequest(w http.ResponseWriter, r *http.Request, idStr string) {
	var request models.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid follow request ID", http.StatusBadRequest)
		return
	}

	if err := services.UpdateFollowRequest(id, request.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteFollowRequest deletes a follow request by ID
func DeleteFollowRequest(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid follow request ID", http.StatusBadRequest)
		return
	}

	if err := services.DeleteFollowRequest(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
