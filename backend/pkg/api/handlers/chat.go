package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	// Check if userID is set correctly
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var message models.Chat
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	message.SenderID = userID

	if message.CreatedAt.IsZero() {
		message.CreatedAt = time.Now()
	}

	log.Printf("Processed message for storage: %+v", message)

	if err := services.SendMessage(message); err != nil {
		http.Error(w, "Failed to send message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Message sent successfully",
	})
}

// GetMessages handles GET requests to retrieve messages for a specific recipient or group
func GetMessages(w http.ResponseWriter, r *http.Request, recipientIDStr string, groupIDStr string) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	recipientID, err := strconv.Atoi(recipientIDStr)
	if err != nil {
		http.Error(w, "Invalid recipient ID", http.StatusBadRequest)
		return
	}

	var groupID int
	if groupIDStr != "" {
		groupID, err = strconv.Atoi(groupIDStr)
		if err != nil {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}
	}

	messages, err := services.GetMessages(userID, recipientID, groupID)
	if err != nil {
		http.Error(w, "Failed to retrieve messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(messages) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Failed to encode messages: "+err.Error(), http.StatusInternalServerError)
	}
}
