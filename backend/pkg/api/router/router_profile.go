package router

import (
	"Social/pkg/api/handlers"
	"net/http"
	"strings"
)

func HandleProfileRoutes(w http.ResponseWriter, r *http.Request) {
	// Extract userID from URL path
	userIDStr := strings.TrimPrefix(r.URL.Path, "/profile/")

	// Check if userID is not empty
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		handlers.GetProfile(w, r, userIDStr)
	case "PUT":
		handlers.UpdateProfile(w, r, userIDStr)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
