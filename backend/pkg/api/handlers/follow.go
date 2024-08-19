package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"encoding/json"
	"net/http"
	"strconv"
)

// CreateFollowRequest handles POST requests to create a follow request
func CreateFollowRequest(w http.ResponseWriter, r *http.Request) {
	var request models.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
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
		http.Error(w, "Failed to create follow request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Follow request created successfully",
	})
}

// GetFollowRequest handles GET requests to retrieve a specific follow request
func GetFollowRequest(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid follow request ID", http.StatusBadRequest)
		return
	}

	followRequest, err := services.GetFollowRequest(id)
	if err != nil {
		http.Error(w, "Follow request not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(followRequest); err != nil {
		http.Error(w, "Failed to encode follow request: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateFollowRequest handles PUT requests to update a follow request
func UpdateFollowRequest(w http.ResponseWriter, r *http.Request, idStr string) {
	var request models.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid follow request ID", http.StatusBadRequest)
		return
	}

	if err := services.UpdateFollowRequest(id, request.Status); err != nil {
		http.Error(w, "Failed to update follow request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Follow request updated successfully",
	})
}

// DeleteFollowRequest handles DELETE requests to delete a follow request
func DeleteFollowRequest(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid follow request ID", http.StatusBadRequest)
		return
	}

	if err := services.DeleteFollowRequest(id); err != nil {
		http.Error(w, "Failed to delete follow request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Follow request deleted successfully",
	})
}

// AcceptFollowRequest handles POST requests to accept a follow request
func AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid follow request ID", http.StatusBadRequest)
		return
	}

	// Update the follow request status to accepted
	if err := services.UpdateFollowRequest(id, models.FollowRequestAccepted); err != nil {
		http.Error(w, "Failed to accept follow request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the follow request
	followRequest, err := services.GetFollowRequest(id)
	if err != nil {
		http.Error(w, "Failed to retrieve follow request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Add recipient to sender's followers list
	if err := services.AddFollower(followRequest.SenderID, followRequest.RecipientID); err != nil {
		http.Error(w, "Failed to add follower: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Follow request accepted and follower added successfully",
	})
}

// RejectFollowRequest handles POST requests to reject a follow request
func RejectFollowRequest(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid follow request ID", http.StatusBadRequest)
		return
	}

	// Update the follow request status to rejected
	if err := services.UpdateFollowRequest(id, models.FollowRequestRejected); err != nil {
		http.Error(w, "Failed to reject follow request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Follow request rejected successfully",
	})
}
