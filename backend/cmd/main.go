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
