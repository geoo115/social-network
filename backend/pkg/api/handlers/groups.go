package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// CreateGroup handles POST requests to create a new group
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	groupID, err := services.CreateGroup(group)
	if err != nil {
		http.Error(w, "Failed to create group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Group created successfully",
		"group_id": groupID,
	})
}

// GetGroup handles GET requests to retrieve a specific group
func GetGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/groups/"))
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	group, err := services.GetGroup(groupID)
	if err != nil {
		http.Error(w, "Group not found: "+err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(group)
}

// ListGroups handles GET requests to list groups with optional filtering
func ListGroups(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")
	searchTerm := r.URL.Query().Get("search")

	offsetInt, _ := strconv.Atoi(offset)
	limitInt, _ := strconv.Atoi(limit)

	groups, err := services.ListGroups(offsetInt, limitInt, searchTerm)
	if err != nil {
		http.Error(w, "Failed to retrieve groups: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(groups)
}

// InviteToGroup handles POST requests to invite a user to a group
func InviteToGroup(w http.ResponseWriter, r *http.Request) {
	var invitation models.GroupInvitation
	err := json.NewDecoder(r.Body).Decode(&invitation)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = services.InviteToGroup(invitation)
	if err != nil {
		http.Error(w, "Failed to invite user to group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User invited to group successfully",
	})
}

// CreateGroupRequest handles POST requests to create a group request
func CreateGroupRequest(w http.ResponseWriter, r *http.Request) {
	var request models.GroupRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = services.CreateGroupRequest(request)
	if err != nil {
		http.Error(w, "Failed to create group request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Group request created successfully",
	})
}

// CreateGroupEvent handles POST requests to create a group event
func CreateGroupEvent(w http.ResponseWriter, r *http.Request) {
	var event models.GroupEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = services.CreateGroupEvent(event)
	if err != nil {
		http.Error(w, "Failed to create group event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Group event created successfully",
	})
}

// RSVPEvent handles POST requests to RSVP to a group event
func RSVPEvent(w http.ResponseWriter, r *http.Request) {
	var rsvp models.EventRSVP
	err := json.NewDecoder(r.Body).Decode(&rsvp)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = services.RSVPEvent(rsvp)
	if err != nil {
		http.Error(w, "Failed to RSVP to event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "RSVP to event successful",
	})
}

// GetGroupEvent handles GET requests to retrieve a specific group event
func GetGroupEvent(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/groups/"), "/")
	if len(pathSegments) < 3 || pathSegments[1] != "events" {
		http.Error(w, "Invalid event request", http.StatusBadRequest)
		return
	}

	groupID, err := strconv.Atoi(pathSegments[0])
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.Atoi(pathSegments[2])
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	event, err := services.GetGroupEvent(groupID, eventID)
	if err != nil {
		http.Error(w, "Event not found: "+err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(event)
}

// JoinGroup handles POST requests to join a group
func JoinGroup(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/groups/")
	pathSegments := strings.Split(path, "/")

	if len(pathSegments) < 2 || pathSegments[1] != "join" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	groupID, err := strconv.Atoi(pathSegments[0])
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		UserID int `json:"user_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = services.JoinGroup(groupID, requestBody.UserID)
	if err != nil {
		http.Error(w, "Failed to join group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully joined the group",
	})
}

// LeaveGroup handles POST requests to leave a group
func LeaveGroup(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/groups/")
	pathSegments := strings.Split(path, "/")

	if len(pathSegments) < 2 || pathSegments[1] != "leave" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	groupID, err := strconv.Atoi(pathSegments[0])
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		UserID int `json:"user_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = services.LeaveGroup(groupID, requestBody.UserID)
	if err != nil {
		http.Error(w, "Failed to leave group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully left the group",
	})
}

// RespondToInvitation handles POST requests to respond to a group invitation
func RespondToInvitation(w http.ResponseWriter, r *http.Request, invitationID int) {
	var response struct {
		Status string `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = services.RespondToInvitation(invitationID, response.Status)
	if err != nil {
		http.Error(w, "Failed to respond to invitation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Invitation response recorded successfully",
	})
}

// RespondToGroupRequest handles POST requests to respond to a group request
func RespondToGroupRequest(w http.ResponseWriter, r *http.Request, requestID int) {
	var response struct {
		Status string `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = services.RespondToGroupRequest(requestID, response.Status)
	if err != nil {
		http.Error(w, "Failed to respond to group request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Group request response recorded successfully",
	})
}
