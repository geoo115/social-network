package router

import (
	"Social/pkg/api/handlers"
	"net/http"
	"strings"
)

func HandlePostRoutes(w http.ResponseWriter, r *http.Request) {
	// Extract postID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/post/")
	parts := strings.SplitN(path, "/", 2)
	postIDStr := ""

	if len(parts) > 0 {
		postIDStr = parts[0]
	}

	switch r.Method {
	case "GET":
		if postIDStr != "" {
			handlers.GetPost(w, r, postIDStr)
		} else {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
		}
	case "PUT":
		if postIDStr != "" {
			handlers.UpdatePost(w, r, postIDStr)
		} else {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
		}
	case "DELETE":
		if postIDStr != "" {
			handlers.DeletePost(w, r, postIDStr)
		} else {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
