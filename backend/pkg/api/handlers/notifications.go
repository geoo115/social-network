package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetNotifications handles GET requests to retrieve notifications for the current user
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notifications, err := services.GetNotifications(userID)
	if err != nil {
		log.Printf("Failed to get notifications for user %d: %v", userID, err)
		http.Error(w, "Failed to retrieve notifications", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(notifications); err != nil {
		log.Printf("Failed to encode notifications: %v", err)
		http.Error(w, "Failed to encode notifications", http.StatusInternalServerError)
		return
	}
}

// CreateNotification handles POST requests to create a new notification
func CreateNotification(w http.ResponseWriter, r *http.Request) {
	var notification models.Notification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := services.CreateNotification(notification); err != nil {
		log.Printf("Failed to create notification: %v", err)
		http.Error(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// NotifyFollowRequest sends a notification when a user receives a new follower request
func NotifyFollowRequest(followerID, followedID int) {
	message := "You have a new follower request."
	notification := models.Notification{
		UserID:    followedID,
		Type:      "follow_request",
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
		Details:   fmt.Sprintf("follower_id:%d", followerID),
	}
	if err := services.CreateNotification(notification); err != nil {
		log.Printf("Failed to send follow request notification: %v", err)
	}
}

// NotifyGroupInvite sends a notification when a user is invited to a group
func NotifyGroupInvite(inviterID, invitedID, groupID int) {
	message := "You have been invited to join a group."
	notification := models.Notification{
		UserID:    invitedID,
		Type:      "group_invite",
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
		Details:   fmt.Sprintf("group_id:%d, inviter_id:%d", groupID, inviterID),
	}
	if err := services.CreateNotification(notification); err != nil {
		log.Printf("Failed to send group invite notification: %v", err)
	}
}

// NotifyGroupJoinRequest sends a notification when a user requests to join a group
func NotifyGroupJoinRequest(groupCreatorID, requesterID, groupID int) {
	message := "A user has requested to join your group."
	notification := models.Notification{
		UserID:    groupCreatorID,
		Type:      "group_join_request",
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
		Details:   fmt.Sprintf("group_id:%d, requester_id:%d", groupID, requesterID),
	}
	if err := services.CreateNotification(notification); err != nil {
		log.Printf("Failed to send group join request notification: %v", err)
	}
}

// NotifyEventCreation sends a notification when an event is created in a group
func NotifyEventCreation(groupID, eventID int) {
	message := "An event has been created in your group."
	notification := models.Notification{
		UserID:    groupID, // Assuming this is a group member or related user ID
		Type:      "event_created",
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
		Details:   fmt.Sprintf("event_id:%d", eventID),
	}
	if err := services.CreateNotification(notification); err != nil {
		log.Printf("Failed to send event creation notification: %v", err)
	}
}

// MarkNotificationAsRead handles marking a notification as read
func MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	notificationIDStr := r.URL.Query().Get("id")
	if notificationIDStr == "" {
		http.Error(w, "Notification ID is required", http.StatusBadRequest)
		return
	}

	notificationID, err := strconv.Atoi(notificationIDStr)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	if err := services.MarkNotificationAsRead(notificationID); err != nil {
		log.Printf("Failed to mark notification %d as read: %v", notificationID, err)
		http.Error(w, "Failed to mark notification as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
