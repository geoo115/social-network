package api

import (
	"Social/pkg/api/handlers"
	"Social/pkg/api/middlewares"
	"encoding/json"
	"net/http"
	"strconv"
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

	// Invitations and Requests routes
	mux.Handle("/invitations/", middlewares.SessionAuthMiddleware(http.HandlerFunc(handleInvitationRoutes)))
	mux.Handle("/requests/", middlewares.SessionAuthMiddleware(http.HandlerFunc(handleRequestRoutes)))

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

func handleInvitationRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/invitations/")
	pathSegments := strings.Split(path, "/")

	if len(pathSegments) < 2 || pathSegments[1] != "response" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	invitationID, err := strconv.Atoi(pathSegments[0])
	if err != nil {
		http.Error(w, "Invalid invitation ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		handlers.RespondToInvitation(w, r, invitationID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleRequestRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/requests/")
	pathSegments := strings.Split(path, "/")

	if len(pathSegments) < 2 || pathSegments[1] != "response" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	requestID, err := strconv.Atoi(pathSegments[0])
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		handlers.RespondToGroupRequest(w, r, requestID)
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
	// Ensure that the path is exactly /notifications
	if r.URL.Path != "/notifications" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handlers.GetNotifications(w, r)
	case http.MethodPost:
		handlers.CreateNotification(w, r)
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
