package router

import (
	"Social/pkg/api/handlers"
	"net/http"
)

func HandleNotificationRoutes(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/notifications" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handlers.GetNotifications(w, r)
	case http.MethodPost:
		handlers.CreateNotification(w, r)
	case http.MethodPut:
		handlers.MarkNotificationAsRead(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
