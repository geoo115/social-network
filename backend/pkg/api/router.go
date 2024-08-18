package api

import (
	"Social/pkg/api/handlers"
	"Social/pkg/api/middlewares"
	"Social/pkg/api/router"
	"net/http"
)

func InitializeRoutes(mux *http.ServeMux) {
	// Authentication routes
	mux.Handle("/register", http.HandlerFunc(handlers.Register))
	mux.Handle("/login", http.HandlerFunc(handlers.Login))

	// Profile routes (protected routes)
	mux.Handle("/profile/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleProfileRoutes)))

	// Posts routes
	mux.Handle("/post", middlewares.SessionAuthMiddleware(http.HandlerFunc(handlers.CreatePost)))
	mux.Handle("/post/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandlePostRoutes)))

	// Like/Dislike routes
	mux.Handle("/posts/like", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleLikeDislikeRoutes)))
	mux.Handle("/posts/dislike", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleLikeDislikeRoutes)))

	// Comments routes
	mux.Handle("/comments/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleCommentRoutes)))

	// Group routes
	mux.Handle("/groups/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleGroupRoutes)))

	// Invitations and Requests routes
	mux.Handle("/invitations/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleInvitationRoutes)))
	mux.Handle("/requests/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleRequestRoutes)))

	// Chat routes
	mux.Handle("/chats/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleChatRoutes)))

	// Notification routes
	mux.Handle("/notifications", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleNotificationRoutes)))

	// Follow Request routes
	mux.Handle("/follow-requests/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleFollowRequestRoutes)))

	// Follow Request Acceptance
	mux.Handle("/follow-requests/accept", middlewares.SessionAuthMiddleware(http.HandlerFunc(handlers.AcceptFollowRequest)))

	// Follow Request Rejection
	mux.Handle("/follow-requests/reject", middlewares.SessionAuthMiddleware(http.HandlerFunc(handlers.RejectFollowRequest)))

	// Fallback for not found routes
	mux.Handle("/", http.NotFoundHandler())
}
