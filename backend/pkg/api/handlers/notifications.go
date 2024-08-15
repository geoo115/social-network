package handlers

import (
	"Social/pkg/services"
	"encoding/json"
	"net/http"
)

// GetNotifications handles GET requests to retrieve user notifications
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

	if err := json.NewEncoder(w).Encode(notifications); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
