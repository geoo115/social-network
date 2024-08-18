package router

import (
	"Social/pkg/api/handlers"
	"encoding/json"
	"net/http"
	"strings"
)

func HandleLikeDislikeRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req struct {
			PostID int `json:"post_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.PostID == 0 {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/posts/like") {
			handlers.LikePost(w, r, req.PostID)
		} else if strings.HasPrefix(r.URL.Path, "/posts/dislike") {
			handlers.DislikePost(w, r, req.PostID)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
