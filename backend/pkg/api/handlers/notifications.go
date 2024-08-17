package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notifications, err := services.GetNotifications(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the notifications
	fmt.Printf("Retrieved notifications: %+v\n", notifications)

	if err := json.NewEncoder(w).Encode(notifications); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateNotification(w http.ResponseWriter, r *http.Request) {
	var notification models.Notification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := services.CreateNotification(notification); err != nil {
		http.Error(w, "Failed to create notification: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
