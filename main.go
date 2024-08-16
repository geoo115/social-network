package main

import (
	"log"
	"net/http"
	"os"

	"Social/pkg/api"
	"Social/pkg/db"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize SQLite connection
	err = db.Initialize()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Initialize the routes from the api/router.go
	mux := http.NewServeMux()
	api.InitializeRoutes(mux)

	// Get the port from the environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	// Start the server
	log.Printf("Server starting on port %s", port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
package api

import (
	"Social/pkg/api/handlers"
	"Social/pkg/api/middlewares"
	"encoding/json"
	"net/http"
	"strings"
)

func InitializeRoutes(mux *http.ServeMux) {
	// Authentication routes
	mux.Handle("/register", http.HandlerFunc(handlers.Register))
	mux.Handle("/login", http.HandlerFunc(handlers.Login))

	// Profile routes (protected routes)
	mux.Handle("/profile/", middlewares.SessionAuthMiddleware(http.HandlerFunc(handleProfileRoutes)))

	// Posts routes
	mux.Handle("/post", middlewares.SessionAuthMiddleware(http.HandlerFunc(handlers.CreatePost)))
	mux.Handle("/post/", middlewares.SessionAuthMiddleware(http.HandlerFunc(handlePostRoutes)))

	// Like/Dislike routes
	mux.Handle("/posts/like", middlewares.SessionAuthMiddleware(http.HandlerFunc(handleLikeDislikeRoutes)))
	mux.Handle("/posts/dislike", middlewares.SessionAuthMiddleware(http.HandlerFunc(handleLikeDislikeRoutes)))

	// Comments routes
	mux.Handle("/comments/", middlewares.SessionAuthMiddleware(http.HandlerFunc(handleCommentRoutes)))

	// Group routes
	mux.Handle("/groups/", middlewares.SessionAuthMiddleware(http.HandlerFunc(handleGroupRoutes)))

	// Chat routes
	mux.Handle("/chats/", middlewares.SessionAuthMiddleware(http.HandlerFunc(handleChatRoutes)))

	// Notification routes
	mux.Handle("/notifications", middlewares.SessionAuthMiddleware(http.HandlerFunc(handleNotificationRoutes)))

	// Follow Request routes
	mux.Handle("/follow-requests/", middlewares.SessionAuthMiddleware(http.HandlerFunc(handleFollowRequestRoutes)))

	// Fallback for not found routes
	mux.Handle("/", http.NotFoundHandler())
}

func handleProfileRoutes(w http.ResponseWriter, r *http.Request) {
	// Extract userID from URL path
	userIDStr := r.URL.Path[len("/profile/"):]

	switch r.Method {
	case "GET":
		handlers.GetProfile(w, r, userIDStr)
	case "PUT":
		handlers.UpdateProfile(w, r, userIDStr)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePostRoutes(w http.ResponseWriter, r *http.Request) {
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

func handleLikeDislikeRoutes(w http.ResponseWriter, r *http.Request) {
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

func handleCommentRoutes(w http.ResponseWriter, r *http.Request) {
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

func handleGroupRoutes(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method
	segments := strings.Split(strings.TrimPrefix(path, "/"), "/")

	// Check for empty path
	if len(segments) == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if segments[0] != "groups" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	switch method {
	case http.MethodPost:
		if len(segments) == 2 {
			switch segments[1] {
			case "events":
				handlers.CreateGroupEvent(w, r)
			case "join":
				handlers.JoinGroup(w, r)
			case "leave":
				handlers.LeaveGroup(w, r)
			default:
				handlers.CreateGroup(w, r)
			}
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}

	case http.MethodGet:
		if len(segments) == 2 {
			handlers.GetGroup(w, r)
		} else if len(segments) == 3 && segments[2] == "events" {
			handlers.GetGroupEvent(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}

	case http.MethodPut:
		if len(segments) == 3 && segments[2] == "events" {
			handlers.UpdateGroupEvent(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}

	case http.MethodDelete:
		if len(segments) == 3 && segments[2] == "events" {
			handlers.DeleteGroupEvent(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleChatRoutes(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	segments := strings.Split(strings.TrimPrefix(path, "/chats/"), "/")

	if len(segments) < 2 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	action := segments[0]
	recipientIDStr := segments[1]
	groupIDStr := ""

	if len(segments) > 2 {
		groupIDStr = segments[2]
	}

	switch r.Method {
	case "POST":
		if action == "send" {
			handlers.SendMessage(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	case "GET":
		if action == "messages" {
			handlers.GetMessages(w, r, recipientIDStr, groupIDStr)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleNotificationRoutes(w http.ResponseWriter, r *http.Request) {
	// Ensure that the path is exactly /notifications for GET requests
	if r.URL.Path != "/notifications" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		handlers.GetNotifications(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleFollowRequestRoutes(w http.ResponseWriter, r *http.Request) {
	// Ensure that the path starts with /follow-requests/
	if !strings.HasPrefix(r.URL.Path, "/follow-requests/") {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Extract the ID from the URL path, if present
	id := strings.TrimPrefix(r.URL.Path, "/follow-requests/")
	id = strings.Trim(id, "/")

	switch r.Method {
	case "POST":
		handlers.CreateFollowRequest(w, r)
	case "GET":
		if id != "" {
			handlers.GetFollowRequest(w, r, id)
		} else {
			http.Error(w, "ID required", http.StatusBadRequest)
		}
	case "PUT":
		if id != "" {
			handlers.UpdateFollowRequest(w, r, id)
		} else {
			http.Error(w, "ID required", http.StatusBadRequest)
		}
	case "DELETE":
		if id != "" {
			handlers.DeleteFollowRequest(w, r, id)
		} else {
			http.Error(w, "ID required", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
package handlers

import (
	"encoding/json"
	"net/http"

	"Social/pkg/api/middlewares"
	"Social/pkg/models"
	"Social/pkg/services"
)

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the service to register the user
	err := services.RegisterUser(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the service to authenticate the user
	user, err := services.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Generate a new session ID
	sessionID, err := middlewares.GenerateSessionID(user.ID)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	// Set the session cookie
	middlewares.SetSessionCookie(w, sessionID)

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}
package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"encoding/json"
	"net/http"
	"strconv"
)

// SendMessage handles POST requests to send a message
func SendMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Chat
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	message.SenderID = userID

	if err := services.SendMessage(message); err != nil {
		http.Error(w, "Failed to send message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetMessages handles GET requests to retrieve messages
func GetMessages(w http.ResponseWriter, r *http.Request, recipientIDStr string, groupIDStr string) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	recipientID, err := strconv.Atoi(recipientIDStr)
	if err != nil {
		http.Error(w, "Invalid recipient ID", http.StatusBadRequest)
		return
	}

	var groupID int
	if groupIDStr != "" {
		groupID, err = strconv.Atoi(groupIDStr)
		if err != nil {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}
	}

	messages, err := services.GetMessages(userID, recipientID, groupID)
	if err != nil {
		http.Error(w, "Failed to retrieve messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(messages) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Failed to encode messages: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"encoding/json"
	"net/http"
	"strconv"
)

// CreateComment handles POST requests to create a comment
func CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	comment.UserID = userID

	if err := services.CreateComment(comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetComment handles GET requests to retrieve a specific comment
func GetComment(w http.ResponseWriter, r *http.Request, commentIDStr string) {
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	comment, err := services.GetComment(commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateComment handles PUT requests to update a comment
func UpdateComment(w http.ResponseWriter, r *http.Request, commentIDStr string) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	if err := services.UpdateComment(commentID, comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteComment handles DELETE requests to delete a comment
func DeleteComment(w http.ResponseWriter, r *http.Request, commentIDStr string) {
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	if err := services.DeleteComment(commentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"encoding/json"
	"net/http"
	"strconv"
)

// CreateFollowRequest handles the creation of a new follow request
func CreateFollowRequest(w http.ResponseWriter, r *http.Request) {
	var request models.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	request.SenderID = userID
	request.Status = models.FollowRequestPending

	if err := services.CreateFollowRequest(request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetFollowRequest retrieves a follow request by ID
func GetFollowRequest(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid follow request ID", http.StatusBadRequest)
		return
	}

	followRequest, err := services.GetFollowRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(followRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateFollowRequest updates the status of an existing follow request
func UpdateFollowRequest(w http.ResponseWriter, r *http.Request, idStr string) {
	var request models.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid follow request ID", http.StatusBadRequest)
		return
	}

	if err := services.UpdateFollowRequest(id, request.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteFollowRequest deletes a follow request by ID
func DeleteFollowRequest(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid follow request ID", http.StatusBadRequest)
		return
	}

	if err := services.DeleteFollowRequest(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
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
package handlers

import (
	"Social/pkg/services"
	"encoding/json"
	"net/http"
)

// GetNotifications handles GET requests to retrieve user notifications
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notifications, err := services.GetNotifications(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(notifications); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
package handlers

import (
	"Social/pkg/models"
	"Social/pkg/services"
	"encoding/json"
	"net/http"
	"strconv"
)

// CreatePost handles POST requests to create a new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Assume userID is set in the context (you should set it in your middleware)
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	post.UserID = userID

	err := services.CreatePost(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"Social/pkg/models"
	"Social/pkg/services"
)

// GetProfile retrieves a user's profile by ID
func GetProfile(w http.ResponseWriter, r *http.Request, userIDStr string) {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	profile, err := services.GetProfile(userID)
	if err != nil {
		if err.Error() == "profile not found" {
			http.Error(w, "Profile not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve profile", http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(profile); err != nil {
		http.Error(w, "Failed to encode profile", http.StatusInternalServerError)
	}
}

// UpdateProfile updates a user's profile information
func UpdateProfile(w http.ResponseWriter, r *http.Request, userIDStr string) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = services.UpdateProfile(userID, user)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
package middlewares

import (
	"context"
	"net/http"
	"time"
)

// SessionAuthMiddleware checks for a valid session ID in the cookie
func SessionAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Missing session ID", http.StatusUnauthorized)
			return
		}

		userID, err := RetrieveSession(sessionID.Value)
		if err != nil {
			http.Error(w, "Invalid or expired session", http.StatusUnauthorized)
			return
		}

		// Add user ID to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)
		r = r.WithContext(ctx)

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}

// SetSessionCookie sets the session ID in a cookie
func SetSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, // Prevent JavaScript access
		Secure:   true, // Set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
	})
}
package middlewares

import (
	"Social/pkg/db"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"
)

// GenerateSessionID generates a new session ID
func GenerateSessionID(userID int) (string, error) {
	sessionID := generateRandomString(32)
	expiresAt := time.Now().Add(24 * time.Hour) // Sessions expire after 24 hours

	_, err := db.DB.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", sessionID, userID, expiresAt)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// RetrieveSession retrieves session data by session ID
func RetrieveSession(sessionID string) (int, error) {
	var userID int
	var expiresAt time.Time

	err := db.DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_id = ?", sessionID).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err
	}

	if time.Now().After(expiresAt) {
		return 0, sql.ErrNoRows // Session expired
	}

	return userID, nil
}

// DeleteSession deletes a session by session ID
func DeleteSession(sessionID string) error {
	_, err := db.DB.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}

// generateRandomString generates a random string of the given length
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

var DB *sql.DB

func Initialize() error {
	var err error
	DB, err = sql.Open("sqlite3", "./socialNetwork.db")
	if err != nil {
		return err
	}

	err = applyMigrations(DB)
	if err != nil {
		return err
	}

	return nil
}

func applyMigrations(db *sql.DB) error {
	tables := []struct {
		name   string
		create string
	}{
		{
			name: "users",
			create: `CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				email TEXT UNIQUE NOT NULL,
				password TEXT NOT NULL,
				first_name TEXT NOT NULL,
				last_name TEXT NOT NULL,
				date_of_birth TEXT NOT NULL,
				avatar TEXT,
				nickname TEXT,
				about_me TEXT,
				created_at DATETIME,
				updated_at DATETIME
			)`,
		},
		{
			name: "posts",
			create: `CREATE TABLE IF NOT EXISTS posts (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				content TEXT NOT NULL,
				image TEXT,
				privacy TEXT NOT NULL,
				created_at DATETIME,
				updated_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id)
			)`,
		},
		{
			name: "groups",
			create: `CREATE TABLE IF NOT EXISTS groups (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				creator_id INTEGER NOT NULL,
				title TEXT NOT NULL,
				description TEXT,
				created_at DATETIME,
				updated_at DATETIME,
				FOREIGN KEY (creator_id) REFERENCES users(id)
			)`,
		},
		{
			name: "chats",
			create: `CREATE TABLE IF NOT EXISTS chats (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				sender_id INTEGER NOT NULL,
				recipient_id INTEGER,
				group_id INTEGER,
				message TEXT NOT NULL,
				is_group BOOLEAN NOT NULL,
				created_at DATETIME,
				FOREIGN KEY (sender_id) REFERENCES users(id),
				FOREIGN KEY (recipient_id) REFERENCES users(id),
				FOREIGN KEY (group_id) REFERENCES groups(id)
			)`,
		},
		{
			name: "notifications",
			create: `CREATE TABLE IF NOT EXISTS notifications (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				message TEXT NOT NULL,
				is_read BOOLEAN NOT NULL,
				created_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id)
			)`,
		},
		{
			name: "follow_requests",
			create: `CREATE TABLE IF NOT EXISTS follow_requests (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				sender_id INTEGER NOT NULL,
				recipient_id INTEGER NOT NULL,
				status TEXT NOT NULL,
				created_at DATETIME,
				FOREIGN KEY (sender_id) REFERENCES users(id),
				FOREIGN KEY (recipient_id) REFERENCES users(id)
			)`,
		},
		{
			name: "group_events",
			create: `CREATE TABLE IF NOT EXISTS group_events (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				group_id INTEGER NOT NULL,
				title TEXT NOT NULL,
				description TEXT,
				day_time DATETIME NOT NULL,
				created_at DATETIME,
				updated_at DATETIME,
				FOREIGN KEY (group_id) REFERENCES groups(id)
			)`,
		},
		{
			name: "group_memberships",
			create: `CREATE TABLE IF NOT EXISTS group_memberships (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				group_id INTEGER NOT NULL,
				joined_at DATETIME,
				left_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (group_id) REFERENCES groups(id)
			)`,
		},
		{
			name: "likes",
			create: `CREATE TABLE IF NOT EXISTS likes (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				post_id INTEGER NOT NULL,
				created_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (post_id) REFERENCES posts(id)
			)`,
		},
		{
			name: "dislikes",
			create: `CREATE TABLE IF NOT EXISTS dislikes (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				post_id INTEGER NOT NULL,
				created_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (post_id) REFERENCES posts(id)
			)`,
		},
		{
			name: "sessions",
			create: `CREATE TABLE IF NOT EXISTS sessions (
				session_id TEXT PRIMARY KEY,
				user_id INTEGER NOT NULL,
				expires_at DATETIME NOT NULL,
				FOREIGN KEY (user_id) REFERENCES users(id)
			)`,
		},
		{
			name: "comments",
			create: `CREATE TABLE IF NOT EXISTS comments (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				post_id INTEGER NOT NULL,
				content TEXT NOT NULL,
				created_at DATETIME,
				updated_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (post_id) REFERENCES posts(id)
			)`,
		},
	}

	for _, table := range tables {
		if _, err := db.Exec(table.create); err != nil {
			return err
		}
	}

	return nil
}
package models

import "time"

const (
	PrivacyPublic        = "public"
	PrivacyPrivate       = "private"
	PrivacyAlmostPrivate = "almost_private"

	FollowRequestPending  = "pending"
	FollowRequestAccepted = "accepted"
	FollowRequestRejected = "rejected"
)

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Avatar      string `json:"avatar,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	AboutMe     string `json:"about_me,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth string    `json:"date_of_birth"`
	Avatar      string    `json:"avatar,omitempty"`
	Nickname    string    `json:"nickname,omitempty"`
	AboutMe     string    `json:"about_me,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	Image     string    `json:"image,omitempty"`
	Privacy   string    `json:"privacy"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Group struct {
	ID          int       `json:"id"`
	CreatorID   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Chat struct {
	ID          int       `json:"id"`
	SenderID    int       `json:"sender_id"`
	RecipientID int       `json:"recipient_id"`
	GroupID     int       `json:"group_id,omitempty"`
	Message     string    `json:"message"`
	IsGroup     bool      `json:"is_group"`
	CreatedAt   time.Time `json:"created_at"`
}

type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type FollowRequest struct {
	ID          int       `json:"id"`
	SenderID    int       `json:"sender_id"`
	RecipientID int       `json:"recipient_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupEvent struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DayTime     time.Time `json:"day_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GroupMembership struct {
	ID       int        `json:"id"`
	UserID   int        `json:"user_id"`
	GroupID  int        `json:"group_id"`
	JoinedAt time.Time  `json:"joined_at"`
	LeftAt   *time.Time `json:"left_at,omitempty"`
}

type Like struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Dislike struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func RegisterUser(user models.RegisterRequest) error {
	var existingUser models.User
	err := db.DB.QueryRow("SELECT id FROM users WHERE email = ?", user.Email).Scan(&existingUser.ID)
	if err == nil {
		return fmt.Errorf("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	dateOfBirth, err := time.Parse("2006-01-02", user.DateOfBirth)
	if err != nil {
		return fmt.Errorf("invalid date of birth format: %w", err)
	}

	_, err = db.DB.Exec(`INSERT INTO users (email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Email,hashedPassword,user.FirstName,user.LastName,dateOfBirth,user.Avatar,user.Nickname,user.AboutMe,time.Now(),time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}


func AuthenticateUser(email, password string) (models.User, error) {
	var user models.User
	row := db.DB.QueryRow("SELECT id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, created_at, updated_at FROM users WHERE email = ?", email)
	err := row.Scan(
		&user.ID,&user.Email,&user.Password,&user.FirstName,&user.LastName,&user.DateOfBirth,&user.Avatar,&user.Nickname,&user.AboutMe,&user.CreatedAt,&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found: %w", err)
		}
		return user, fmt.Errorf("error retrieving user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, fmt.Errorf("invalid password: %w", err)
	}

	return user, nil
}

// SendMessage sends a message to a user or group
func SendMessage(message models.Chat) error {
	query := `
		INSERT INTO chats (sender_id, recipient_id, group_id, message, is_group, created_at) 
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := db.DB.Exec(query, message.SenderID, message.RecipientID, message.GroupID, message.Message, message.IsGroup, time.Now())
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// GetMessages retrieves messages between two users or in a group
func GetMessages(userID, recipientID int, groupID int) ([]models.Chat, error) {
	var messages []models.Chat

	query := `
		SELECT id, sender_id, recipient_id, group_id, message, is_group, created_at
		FROM chats
		WHERE (sender_id = ? AND recipient_id = ?)
		   OR (sender_id = ? AND recipient_id = ?)
		   OR (group_id = ? AND is_group = ?)
		ORDER BY created_at`

	rows, err := db.DB.Query(query, userID, recipientID, recipientID, userID, groupID, true)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve messages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var msg models.Chat
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.RecipientID, &msg.GroupID, &msg.Message, &msg.IsGroup, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating rows: %w", err)
	}

	return messages, nil
}

func CreateComment(comment models.Comment) error {
	query := `
		INSERT INTO comments (user_id, post_id, content, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)`

	_, err := db.DB.Exec(query, comment.UserID, comment.PostID, comment.Content, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

func GetComment(commentID int) (models.Comment, error) {
	var comment models.Comment

	query := `
		SELECT id, user_id, post_id, content, created_at, updated_at
		FROM comments
		WHERE id = ?`

	row := db.DB.QueryRow(query, commentID)
	err := row.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
	if err == sql.ErrNoRows {
		return comment, fmt.Errorf("comment not found")
	} else if err != nil {
		return comment, fmt.Errorf("failed to retrieve comment: %w", err)
	}

	return comment, nil
}


func UpdateComment(commentID int, updatedComment models.Comment) error {
	query := `
		UPDATE comments
		SET content = ?, updated_at = ?
		WHERE id = ?`

	_, err := db.DB.Exec(query, updatedComment.Content, time.Now(), commentID)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}


func DeleteComment(commentID int) error {
	query := `
		DELETE FROM comments
		WHERE id = ?`

	_, err := db.DB.Exec(query, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

func CreateFollowRequest(request models.FollowRequest) error {
	request.CreatedAt = time.Now()

	_, err := db.DB.Exec(`INSERT INTO follow_requests (sender_id, recipient_id, status, created_at)
		VALUES (?, ?, ?, ?)`, request.SenderID, request.RecipientID, request.Status, request.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create follow request: %w", err)
	}
	return nil
}

func GetFollowRequest(id int) (models.FollowRequest, error) {
	row := db.DB.QueryRow(`SELECT id, sender_id, recipient_id, status, created_at
		FROM follow_requests WHERE id = ?`, id)

	var request models.FollowRequest
	if err := row.Scan(&request.ID, &request.SenderID, &request.RecipientID, &request.Status, &request.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return request, fmt.Errorf("follow request not found")
		}
		return request, fmt.Errorf("failed to get follow request: %w", err)
	}
	return request, nil
}

func UpdateFollowRequest(id int, status string) error {
	_, err := db.DB.Exec(`UPDATE follow_requests SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("failed to update follow request: %w", err)
	}
	return nil
}


func DeleteFollowRequest(id int) error {
	_, err := db.DB.Exec(`DELETE FROM follow_requests WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete follow request: %w", err)
	}
	return nil
}

// LikePost adds a like to a post
func LikePost(userID, postID int) error {
	like := models.Like{
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now(),
	}

	_, err := db.DB.Exec(`INSERT INTO likes (user_id, post_id, created_at) VALUES (?, ?, ?)`,
		like.UserID, like.PostID, like.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to add like: %w", err)
	}
	return nil
}


func DislikePost(userID, postID int) error {
	dislike := models.Dislike{
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now(),
	}

	_, err := db.DB.Exec(`INSERT INTO dislikes (user_id, post_id, created_at) VALUES (?, ?, ?)`,
		dislike.UserID, dislike.PostID, dislike.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to add dislike: %w", err)
	}
	return nil
}

func GetNotifications(userID int) ([]models.Notification, error) {
	var notifications []models.Notification

	rows, err := db.DB.Query(`SELECT id, user_id, message, is_read, created_at FROM notifications WHERE user_id = ?`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var notification models.Notification
		if err := rows.Scan(&notification.ID, &notification.UserID, &notification.Message, &notification.IsRead, &notification.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over notifications: %w", err)
	}

	return notifications, nil
}


func CreatePost(post models.Post) error {
	_, err := db.DB.Exec(`INSERT INTO posts (user_id, content, image, privacy, created_at, updated_at) 
        VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))`,
		post.UserID, post.Content, post.Image, post.Privacy)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

func GetPost(postID int) (models.Post, error) {
	row := db.DB.QueryRow(`SELECT id, user_id, content, image, privacy, created_at, updated_at 
		FROM posts WHERE id = ?`, postID)

	var post models.Post
	err := row.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt, &post.UpdatedAt)
	if err == sql.ErrNoRows {
		return post, fmt.Errorf("post not found")
	} else if err != nil {
		return post, fmt.Errorf("failed to retrieve post: %w", err)
	}

	return post, nil
}

func UpdatePost(postID int, updatedPost models.Post) error {
	_, err := db.DB.Exec(`UPDATE posts SET content = ?, image = ?, privacy = ?, updated_at = ? 
		WHERE id = ?`, updatedPost.Content, updatedPost.Image, updatedPost.Privacy, time.Now(), postID)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

// DeletePost removes a post from the database
func DeletePost(postID int) error {
	_, err := db.DB.Exec(`DELETE FROM posts WHERE id = ?`, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

// GetProfile retrieves a user's profile information
func GetProfile(userID int) (models.User, error) {
	var user models.User
	row := db.DB.QueryRow(`
		SELECT id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, created_at, updated_at
		FROM users WHERE id = ?`, userID)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.Avatar,
		&user.Nickname,
		&user.AboutMe,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("profile not found")
		}
		return user, fmt.Errorf("failed to get profile: %w", err)
	}
	return user, nil
}

// UpdateProfile updates user profile information
func UpdateProfile(userID int, updatedProfile models.User) error {
	_, err := db.DB.Exec(`UPDATE users SET first_name = ?, last_name = ?, date_of_birth = ?, avatar = ?, nickname = ?, about_me = ?, updated_at = ? 
		WHERE id = ?`, updatedProfile.FirstName, updatedProfile.LastName, updatedProfile.DateOfBirth, updatedProfile.Avatar, updatedProfile.Nickname, updatedProfile.AboutMe, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	return nil
}
