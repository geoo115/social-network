package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"Social/pkg/api/middlewares"
	"Social/pkg/models"
	"Social/pkg/services"
	"errors"
	"log"

	"golang.org/x/oauth2"
)

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate request payload
	if err := validateRegisterRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the service to register the user
	if err := services.RegisterUser(req); err != nil {
		if errors.Is(err, services.ErrEmailInUse) {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}
		log.Printf("Error registering user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate request payload
	if err := validateLoginRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the service to authenticate the user
	user, err := services.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		log.Printf("Error authenticating user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Generate a new session ID
	sessionID, err := middlewares.GenerateSessionID(user.ID)
	if err != nil {
		log.Printf("Error generating session ID: %v", err)
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	// Set the session cookie
	middlewares.SetSessionCookie(w, sessionID)

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

// GoogleLogin initiates Google OAuth2 login
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := services.GoogleOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

// GoogleCallback handles Google OAuth2 callback
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := services.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := services.GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// Handle user login/registration using userInfo
	handleOAuthUserLogin(w, userInfo, "google")
}

// FacebookLogin initiates Facebook OAuth2 login
func FacebookLogin(w http.ResponseWriter, r *http.Request) {
	url := services.FacebookOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

// FacebookCallback handles Facebook OAuth2 callback
func FacebookCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := services.FacebookOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := services.FacebookOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// Handle user login/registration using userInfo
	handleOAuthUserLogin(w, userInfo, "facebook")
}

// GitHubLogin initiates GitHub OAuth2 login
func GitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := services.GitHubOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

// GitHubCallback handles GitHub OAuth2 callback
func GitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := services.GitHubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := services.GitHubOAuthConfig.Client(context.Background(), token)
	// Get user info
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// Get email addresses
	resp, err = client.Get("https://api.github.com/user/emails")
	if err != nil {
		http.Error(w, "Failed to get user emails", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var emails []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		http.Error(w, "Failed to decode user emails", http.StatusInternalServerError)
		return
	}

	// Find the primary email address
	var primaryEmail string
	for _, email := range emails {
		if email["primary"].(bool) {
			primaryEmail = email["email"].(string)
			break
		}
	}

	if primaryEmail == "" {
		http.Error(w, "Email not found in user info", http.StatusBadRequest)
		return
	}

	// Handle user login/registration using userInfo
	handleOAuthUserLogin(w, map[string]interface{}{"email": primaryEmail}, "github")
}

// handleOAuthUserLogin processes user info and either logs in or registers the user
func handleOAuthUserLogin(w http.ResponseWriter, userInfo map[string]interface{}, provider string) {
	email, ok := userInfo["email"].(string)
	if !ok {
		http.Error(w, "Email not found in user info", http.StatusBadRequest)
		return
	}

	// Check if user exists, if not create a new user
	user, err := services.FindOrCreateUserByEmail(email, provider)
	if err != nil {
		log.Printf("Error handling user login: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Generate session ID for the authenticated user
	sessionID, err := middlewares.GenerateSessionID(user.ID)
	if err != nil {
		log.Printf("Error generating session ID: %v", err)
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	// Set the session cookie
	middlewares.SetSessionCookie(w, sessionID)

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}
