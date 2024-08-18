package router

import (
	"Social/pkg/api/handlers"
	"net/http"
	"strings"
)

func HandleCommentRoutes(w http.ResponseWriter, r *http.Request) {
	// Extract commentID from URL path for GET, PUT, and DELETE routes
	path := r.URL.Path
	if strings.HasPrefix(path, "/comments/") {
		// Extract commentID from the path
		commentIDStr := strings.TrimPrefix(path, "/comments/")
		commentIDStr = strings.TrimPrefix(commentIDStr, "/")

		switch r.Method {
		case "POST":
			handlers.CreateComment(w, r)
		case "GET":
			if commentIDStr != "" {
				handlers.GetComment(w, r, commentIDStr)
			} else {
				http.Error(w, "Comment ID is required", http.StatusBadRequest)
			}
		case "PUT":
			if commentIDStr != "" {
				handlers.UpdateComment(w, r, commentIDStr)
			} else {
				http.Error(w, "Comment ID is required", http.StatusBadRequest)
			}
		case "DELETE":
			if commentIDStr != "" {
				handlers.DeleteComment(w, r, commentIDStr)
			} else {
				http.Error(w, "Comment ID is required", http.StatusBadRequest)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}
