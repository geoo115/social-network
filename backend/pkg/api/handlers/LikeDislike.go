package handlers

import (
	"net/http"

	"Social/pkg/services"
)

// LikePost handles POST requests to like a post
func LikePost(w http.ResponseWriter, r *http.Request, postID int) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := services.LikePost(userID, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DislikePost handles POST requests to dislike a post
func DislikePost(w http.ResponseWriter, r *http.Request, postID int) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := services.DislikePost(userID, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
