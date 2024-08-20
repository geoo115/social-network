package middlewares

import (
	"context"
	"net/http"
	"time"
)

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

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// SetSessionCookie sets the session ID in a cookie
func SetSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Set to true for production with HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}
