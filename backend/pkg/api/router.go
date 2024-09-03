package api

import (
	"Social/pkg/api/handlers"
	"Social/pkg/api/middlewares"
	"Social/pkg/api/router"
	"net/http"
)

func InitializeRoutes(mux *http.ServeMux) {
	mux.Handle("/register", http.HandlerFunc(handlers.Register))
	mux.Handle("/login", http.HandlerFunc(handlers.Login))

	mux.Handle("/auth/google/login", http.HandlerFunc(handlers.GoogleLogin))
	mux.Handle("/auth/google/callback", http.HandlerFunc(handlers.GoogleCallback))

	mux.Handle("/auth/facebook/login", http.HandlerFunc(handlers.FacebookLogin))
	mux.Handle("/auth/facebook/callback", http.HandlerFunc(handlers.FacebookCallback))

	mux.Handle("/auth/github/login", http.HandlerFunc(handlers.GitHubLogin))
	mux.Handle("/auth/github/callback", http.HandlerFunc(handlers.GitHubCallback))

	mux.Handle("/profile/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleProfileRoutes)))

	mux.Handle("/post", middlewares.SessionAuthMiddleware(http.HandlerFunc(handlers.CreatePost)))
	mux.Handle("/post/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandlePostRoutes)))

	mux.Handle("/posts/like", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleLikeDislikeRoutes)))
	mux.Handle("/posts/dislike", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleLikeDislikeRoutes)))

	mux.Handle("/comments/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleCommentRoutes)))

	mux.Handle("/groups/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleGroupRoutes)))

	mux.Handle("/invitations/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleInvitationRoutes)))
	mux.Handle("/requests/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleRequestRoutes)))

	mux.Handle("/chats/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleChatRoutes)))

	mux.Handle("/notifications", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleNotificationRoutes)))

	mux.Handle("/follow-requests/", middlewares.SessionAuthMiddleware(http.HandlerFunc(router.HandleFollowRequestRoutes)))

	mux.Handle("/follow-requests/accept", middlewares.SessionAuthMiddleware(http.HandlerFunc(handlers.AcceptFollowRequest)))

	mux.Handle("/follow-requests/reject", middlewares.SessionAuthMiddleware(http.HandlerFunc(handlers.RejectFollowRequest)))

	mux.Handle("/", http.NotFoundHandler())
}
