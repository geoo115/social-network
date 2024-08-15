package main

import (
	"log"
	"net/http"
	"os"

	"Social/pkg/api"
	"Social/pkg/db"

	"github.com/gorilla/mux"
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

	// Initialize the router
	router := mux.NewRouter()

	// Initialize the routes from the api/router.go
	api.InitializeRoutes(router)

	// Get the port from the environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	// Start the server
	log.Printf("Server starting on port %s", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
