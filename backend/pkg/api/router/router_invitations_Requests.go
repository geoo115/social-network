package router

import (
	"Social/pkg/api/handlers"
	"net/http"
	"strconv"
	"strings"
)

func HandleInvitationRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/invitations/")
	pathSegments := strings.Split(path, "/")

	if len(pathSegments) < 2 || pathSegments[1] != "response" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	invitationID, err := strconv.Atoi(pathSegments[0])
	if err != nil {
		http.Error(w, "Invalid invitation ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		handlers.RespondToInvitation(w, r, invitationID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleRequestRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/requests/")
	pathSegments := strings.Split(path, "/")

	if len(pathSegments) < 2 || pathSegments[1] != "response" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	requestID, err := strconv.Atoi(pathSegments[0])
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		handlers.RespondToGroupRequest(w, r, requestID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
