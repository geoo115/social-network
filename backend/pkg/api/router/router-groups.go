package router

import (
	"Social/pkg/api/handlers"
	"net/http"
	"strings"
)

func HandleGroupRoutes(w http.ResponseWriter, r *http.Request) {
	// Extract the path after "/groups/"
	path := strings.TrimPrefix(r.URL.Path, "/groups/")
	pathSegments := strings.Split(path, "/")

	if len(pathSegments) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if len(pathSegments) == 2 && pathSegments[1] == "events" {
			handlers.GetGroupEvent(w, r) // Handle GET /groups/{groupID}/events/{eventID}
		} else {
			handlers.GetGroup(w, r) // Handle GET /groups/{groupID}
		}
	case http.MethodPost:
		if len(pathSegments) == 1 {
			handlers.CreateGroup(w, r) // Handle POST /groups/
		} else if len(pathSegments) == 2 && pathSegments[1] == "join" {
			handlers.JoinGroup(w, r) // Handle POST /groups/{groupID}/join
		} else if len(pathSegments) == 2 && pathSegments[1] == "leave" {
			handlers.LeaveGroup(w, r) // Handle POST /groups/{groupID}/leave
		} else if len(pathSegments) == 2 && pathSegments[1] == "events" {
			handlers.CreateGroupEvent(w, r) // Handle POST /groups/{groupID}/events
		} else {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
