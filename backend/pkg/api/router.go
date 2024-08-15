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
