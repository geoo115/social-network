package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	groupID, err := services.CreateGroup(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"group_id": groupID})
}

func GetGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/groups/"))
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	group, err := services.GetGroup(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(group)
}

func ListGroups(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")
	searchTerm := r.URL.Query().Get("search")

	offsetInt, _ := strconv.Atoi(offset)
	limitInt, _ := strconv.Atoi(limit)

	groups, err := services.ListGroups(offsetInt, limitInt, searchTerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(groups)
}

func InviteToGroup(w http.ResponseWriter, r *http.Request) {
	var invitation models.GroupInvitation
	err := json.NewDecoder(r.Body).Decode(&invitation)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = services.InviteToGroup(invitation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func CreateGroupRequest(w http.ResponseWriter, r *http.Request) {
	var request models.GroupRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = services.CreateGroupRequest(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func CreateGroupEvent(w http.ResponseWriter, r *http.Request) {
	var event models.GroupEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = services.CreateGroupEvent(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func RSVPEvent(w http.ResponseWriter, r *http.Request) {
	var rsvp models.EventRSVP
	err := json.NewDecoder(r.Body).Decode(&rsvp)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = services.RSVPEvent(rsvp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

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
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(event)
}

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
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = services.JoinGroup(groupID, requestBody.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = services.LeaveGroup(groupID, requestBody.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RespondToInvitation(w http.ResponseWriter, r *http.Request, invitationID int) {
	var response struct {
		Status string `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = services.RespondToInvitation(invitationID, response.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RespondToGroupRequest(w http.ResponseWriter, r *http.Request, requestID int) {
	var response struct {
		Status string `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = services.RespondToGroupRequest(requestID, response.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
