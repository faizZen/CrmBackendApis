package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	initializers "github.com/Zenithive/it-crm-backend/Initializers"
	"github.com/Zenithive/it-crm-backend/models"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"gorm.io/gorm"
)

// GoogleResponse represents the structure of Google user info
type GoogleResponse struct {
	ID      string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func InitGoogleStore() {
	var store *sessions.CookieStore
	key := "your-session-key" // Replace with a secure key in production
	store = sessions.NewCookieStore([]byte(key))
	store.MaxAge(86400) // 1 day

	// Set up Gothic
	gothic.Store = store

	// Initialize goth with Google provider
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			// "http://localhost:8080/auth/google/callback",
			"http://localhost:8080/auth/google/callback",
			"https://www.googleapis.com/auth/calendar.events.readonly",
			"email", "profile",
		),
	)
}

// VerifyGoogleToken decodes and verifies the Google token
func VerifyGoogleToken(token string) (*GoogleResponse, error) {
	url := "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + token
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("failed to verify Google token")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid Google token")
	}

	var googleData GoogleResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleData); err != nil {
		return nil, errors.New("failed to parse Google response")
	}

	return &googleData, nil
}

// // HandleGoogleLogin handles OAuth login or signup
// func HandleGoogleLogin(googleToken string) (string, string, error) {
// 	googleData, err := VerifyGoogleToken(googleToken)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	var user models.UserDemo
// 	result := initializers.DB.Where("google_id = ? OR email = ?", googleData.ID, googleData.Email).First(&user)

// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		// If user does not exist, create a new user
// 		user = models.UserDemo{
// 			ID:                 uuid.New(),
// 			Email:              googleData.Email,
// 			Name:               googleData.Name,
// 			Password:           "", // No password for Google login
// 			GoogleId:           googleData.ID,
// 			GoogleRefreshToken: googleToken,
// 			// Provider:            googleData.,
// 			BackendRefreshToken: "",
// 			BackendTokenExpiry:  time.Time{},
// 			Role:                "SALES_EXECUTIVE", // Default role
// 		}
// 		// Save new user to DB
// 		if err := initializers.DB.Create(&user).Error; err != nil {
// 			return "", "", fmt.Errorf("failed to create user: %v", err)
// 		}
// 	}
// 	accessToken, refreshToken, err := GenerateTokensDemo(&user, "Google")
// 	if err != nil {
// 		return "", "", err
// 	}
// 	StoreRefreshToken(user.ID.String(), refreshToken)
// 	return accessToken, refreshToken, nil
// }

func OauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Get user from Gothic
	log.Println("OAuth callback reached")
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Println("OAuth error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("User authenticated:", gothUser.Email)

	// Check if user exists by email
	var user models.UserDemo
	result := initializers.DB.Where("email = ?", gothUser.Email).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new user
		user = models.UserDemo{
			ID:                  uuid.New(),
			Email:               gothUser.Email,
			Name:                gothUser.Name,
			Password:            "", // No password for Google login
			GoogleId:            gothUser.UserID,
			GoogleRefreshToken:  gothUser.RefreshToken,
			Provider:            gothUser.Provider,
			BackendRefreshToken: "",
			BackendTokenExpiry:  time.Time{},       //left to set expiration time
			Role:                "SALES_EXECUTIVE", // Default role
		}

		result = initializers.DB.Create(&user)
		if result.Error != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	} else if result.Error != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	log.Println("User :", &user)
	_, refreshToken, err := GenerateTokensDemo(&user, "Google")
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// Store refresh token in DB
	user.BackendRefreshToken = refreshToken
	initializers.DB.Save(&user)

	// Redirect with tokens or return JSON
	// For example, redirect to frontend with tokens as query parameters
	http.Redirect(w, r, fmt.Sprintf("http://localhost:8080/oauth-success?access_token=%s&refresh_token=%s&user_id=%s&auth_provider=%s",
		gothUser.AccessToken, refreshToken, user.ID.String(), user.Provider), http.StatusTemporaryRedirect)
}
