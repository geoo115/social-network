package api

import (
	"net/http"

	"Social/pkg/api/handlers"
	"Social/pkg/api/middlewares"

	"github.com/gorilla/mux"
)

func InitializeRoutes(router *mux.Router) {
	// Authentication routes
	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")

	// Profile routes (protected routes)
	profile := router.PathPrefix("/profile").Subrouter()
	profile.Use(middlewares.AuthMiddleware)
	profile.HandleFunc("/{id}", handlers.GetProfile).Methods("GET")
	profile.HandleFunc("/{id}", handlers.UpdateProfile).Methods("PUT")

	// Posts routes
	posts := router.PathPrefix("/posts").Subrouter()
	posts.Use(middlewares.AuthMiddleware)
	posts.HandleFunc("/", handlers.CreatePost).Methods("POST")
	posts.HandleFunc("/{id}", handlers.GetPost).Methods("GET")
	posts.HandleFunc("/{id}", handlers.UpdatePost).Methods("PUT")
	posts.HandleFunc("/{id}", handlers.DeletePost).Methods("DELETE")

	// Like/Dislike routes
	posts.HandleFunc("/{id}/like", handlers.LikePost).Methods("POST")
	posts.HandleFunc("/{id}/dislike", handlers.DislikePost).Methods("POST")

	// Comments routes
	comments := router.PathPrefix("/comments").Subrouter()
	comments.Use(middlewares.AuthMiddleware)
	comments.HandleFunc("/", handlers.CreateComment).Methods("POST")
	comments.HandleFunc("/{id}", handlers.GetComment).Methods("GET")
	comments.HandleFunc("/{id}", handlers.UpdateComment).Methods("PUT")
	comments.HandleFunc("/{id}", handlers.DeleteComment).Methods("DELETE")

	// Group routes
	groups := router.PathPrefix("/groups").Subrouter()
	groups.Use(middlewares.AuthMiddleware)
	groups.HandleFunc("/", handlers.CreateGroup).Methods("POST")
	groups.HandleFunc("/{id}", handlers.GetGroup).Methods("GET")
	groups.HandleFunc("/{id}/join", handlers.JoinGroup).Methods("POST")
	groups.HandleFunc("/{id}/leave", handlers.LeaveGroup).Methods("POST")
	groups.HandleFunc("/{id}/events", handlers.CreateGroupEvent).Methods("POST")
	groups.HandleFunc("/{id}/events/{eventId}", handlers.GetGroupEvent).Methods("GET")
	groups.HandleFunc("/{id}/events/{eventId}", handlers.UpdateGroupEvent).Methods("PUT")
	groups.HandleFunc("/{id}/events/{eventId}", handlers.DeleteGroupEvent).Methods("DELETE")

	// Chat routes
	chat := router.PathPrefix("/chat").Subrouter()
	chat.Use(middlewares.AuthMiddleware)
	chat.HandleFunc("/{userId}", handlers.SendMessage).Methods("POST")
	chat.HandleFunc("/{userId}/messages", handlers.GetMessages).Methods("GET")

	// Notification routes
	notifications := router.PathPrefix("/notifications").Subrouter()
	notifications.Use(middlewares.AuthMiddleware)
	notifications.HandleFunc("/", handlers.GetNotifications).Methods("GET")

	// Follow Request routes
	followRequests := router.PathPrefix("/follow-requests").Subrouter()
	followRequests.Use(middlewares.AuthMiddleware)
	followRequests.HandleFunc("/", handlers.CreateFollowRequest).Methods("POST")
	followRequests.HandleFunc("/{id}", handlers.GetFollowRequest).Methods("GET")
	followRequests.HandleFunc("/{id}", handlers.UpdateFollowRequest).Methods("PUT")
	followRequests.HandleFunc("/{id}", handlers.DeleteFollowRequest).Methods("DELETE")

	// Fallback for not found routes
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 - Not Found", http.StatusNotFound)
	})
}
