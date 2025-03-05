package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	initializers "github.com/Zenithive/it-crm-backend/Initializers"
	"github.com/Zenithive/it-crm-backend/models"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// GoogleResponse represents the structure of Google user info
type GoogleResponse struct {
	ID      string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
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

// HandleGoogleLogin handles OAuth login or signup
func HandleGoogleLogin(googleToken string) (string, string, error) {
	googleData, err := VerifyGoogleToken(googleToken)
	if err != nil {
		return "", "", err
	}

	var user models.User
	result := initializers.DB.Where("google_id = ? OR email = ?", googleData.ID, googleData.Email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// If user does not exist, create a new user
		user = models.User{
			Email:    googleData.Email,
			Name:     googleData.Name,
			GoogleId: googleData.ID,
			Role:     "user", // Default role
		}

		// Save new user to DB
		if err := initializers.DB.Create(&user).Error; err != nil {
			return "", "", fmt.Errorf("failed to create user: %v", err)
		}
	}

	// Generate tokens
	accessToken, err := GenerateJWT(&user, "google", 1, []byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(user.ID.String())
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// GenerateRefreshToken creates a refresh token and stores it in DB
func GenerateRefreshToken(userID string) (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	expirationTime := time.Now().Add(30 * 24 * time.Hour).Unix() // 30-day expiry

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	// Store refresh token in DB
	refreshTokenRecord := models.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Unix(expirationTime, 0),
	}

	if err := initializers.DB.Create(&refreshTokenRecord).Error; err != nil {
		return "", err
	}

	return refreshToken, nil
}
