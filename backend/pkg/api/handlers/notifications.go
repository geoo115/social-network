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
	services.CreateNotification(notification)
}

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
	services.CreateNotification(notification)
}

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
	services.CreateNotification(notification)
}

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
	services.CreateNotification(notification)
}

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
		http.Error(w, "Failed to mark notification as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
