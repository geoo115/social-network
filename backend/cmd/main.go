package main

import (
	"log"
	"net/http"
	"os"

	"Social/pkg/api"
	"Social/pkg/api/handlers"
	"Social/pkg/api/middlewares"
	"Social/pkg/db"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize SQLite connection
	err = db.Initialize()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Initialize the routes
	mux := http.NewServeMux()
	api.InitializeRoutes(mux)
	corsMux := middlewares.EnableCORS(mux)

	// Start the WebSocket message handling goroutine
	go handlers.HandleMessages()

	// Get the port from the environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	// Start the server
	log.Printf("Server starting on port %s", port)
	err = http.ListenAndServe(":"+port, corsMux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
