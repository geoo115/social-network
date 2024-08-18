package router

import (
	"Social/pkg/api/handlers"
	"net/http"
	"strings"
)

func HandleChatRoutes(w http.ResponseWriter, r *http.Request) {
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
