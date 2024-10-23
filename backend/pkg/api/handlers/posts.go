package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // Limit your max memory size

	// Retrieve the userID from context (set by middleware)
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized: User ID not found", http.StatusUnauthorized)
		return
	}
	log.Printf("Received userID from context: %d", userID)

	// Retrieve form fields
	content := r.FormValue("content")
	privacy := r.FormValue("privacy")

	// Retrieve the file from the form-data
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to retrieve the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Define the path to save the file
	uploadDir := "../../uploads"
	filePath := fmt.Sprintf("%s/%s", uploadDir, handler.Filename)

	// Ensure the upload directory exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
			return
		}
	}

	// Create a file in the defined path
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Copy the file content to the destination file
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return
	}

	// Create a new post object
	post := models.Post{
		UserID:  userID,
		Content: content,
		Image:   filePath, // Store the image path
		Privacy: privacy,
	}

	// Pass the post to the service layer for database insertion
	err = services.CreatePost(post)
	if err != nil {
		log.Printf("Failed to create post: %v", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully"})
}

// GetPost handles GET requests to retrieve a specific post
func GetPost(w http.ResponseWriter, r *http.Request, postIDStr string) {
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := services.GetPost(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, "Failed to encode post", http.StatusInternalServerError)
	}
}

// GetAllPosts handles retrieving all posts
func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	// Retrieve all posts from the service layer
	posts, err := services.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	// Encode posts as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
	}
}

// UpdatePost handles PUT requests to update a post
func UpdatePost(w http.ResponseWriter, r *http.Request, postIDStr string) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	err = services.UpdatePost(postID, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post Updated successfully"})
}

// DeletePost handles DELETE requests to delete a post
func DeletePost(w http.ResponseWriter, r *http.Request, postIDStr string) {
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	err = services.DeletePost(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Delete Post successfully"})
}
