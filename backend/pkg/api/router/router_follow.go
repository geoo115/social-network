package router

import (
	"Social/pkg/api/handlers"
	"log"
	"net/http"
	"strings"
)

func HandleFollowRequestRoutes(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request Method: %s, Request Path: %s", r.Method, r.URL.Path)

	// Ensure the path has the expected prefix
	if !strings.HasPrefix(r.URL.Path, "/follow-requests/") {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Extract the ID, handling the trailing slash
	id := strings.TrimPrefix(r.URL.Path, "/follow-requests/")
	id = strings.Trim(id, "/")

	log.Printf("Parsed ID: %s", id)

	switch r.Method {
	case http.MethodPost:
		handlers.CreateFollowRequest(w, r)
	case http.MethodGet:
		if id != "" {
			handlers.GetFollowRequest(w, r, id)
		} else {
			http.Error(w, "ID required", http.StatusBadRequest)
		}
	case http.MethodPut:
		if id != "" {
			handlers.UpdateFollowRequest(w, r, id)
		} else {
			http.Error(w, "ID required", http.StatusBadRequest)
		}
	case http.MethodDelete:
		if id != "" {
			handlers.DeleteFollowRequest(w, r, id)
		} else {
			http.Error(w, "ID required", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
