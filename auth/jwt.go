package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	initializers "github.com/Zenithive/it-crm-backend/Initializers"
	"github.com/Zenithive/it-crm-backend/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
	// "github.com/golang-jwt/jwt/v5"
)

// var secretKey = []byte(os.Getenv("SECRET_KEY")) // Change this to a secure key
var SecretKey = "abcdefghijklmnopqrstuvwxyz"

// var refreshSecretKey = []byte(os.Getenv("REFRESH_SECRET_KEY"))

var RefreshSecretKey = "This is my RefreshSecretKey"

// Claims structure for JWT
type Claims struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Provider string `json:"provider"`
	jwt.RegisteredClaims
}

func ExtractUserID(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		fmt.Println("Missing token")
		return "", fmt.Errorf("missing token")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method in token")

			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		fmt.Println("Invalid token:", err)
		return "", fmt.Errorf("invalid token: %v", err)
	}

	// Extract user_id from claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Valid token")
		userID, ok := claims["user_id"] // Extract user_id from claims
		if !ok {
			fmt.Println("Invalid user ID in token")
			return "", fmt.Errorf("invalid user ID in token")
		}
		return userID.(string), nil
	}

	fmt.Println("Invalid token claims")
	return "", fmt.Errorf("invalid token claims")
}

// GenerateJWT generates a new token

// StoreRefreshToken saves refresh token in the database
func StoreRefreshToken(userID string, refreshToken string) error {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	var existingToken models.RefreshToken
	result := initializers.DB.Where("user_id = ?", parsedUserID).First(&existingToken)

	if result.Error != nil {
		// If record is not found, create a new entry
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			refreshRecord := models.RefreshToken{
				ID:        uuid.New(),
				UserID:    parsedUserID,
				Token:     refreshToken,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // Expiry (7 days)
			}
			return initializers.DB.Create(&refreshRecord).Error
		}
		// Return error if it's not a "record not found" error
		return result.Error
	}

	// If record exists, update the token and expiration time
	existingToken.Token = refreshToken
	existingToken.CreatedAt = time.Now()
	existingToken.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	return initializers.DB.Save(&existingToken).Error
}

// ValidateRefreshToken checks if the refresh token is valid
func ValidateRefreshToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(RefreshSecretKey), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("refresh token expired")
		}
		return nil, fmt.Errorf("invalid refresh token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token using a valid refresh token
func RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Retrieve user info
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid user ID in refresh token")
	}
	var user models.User
	result := initializers.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return "", errors.New("user not found")
	}

	// Generate a new access token
	accessExpiry, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_TIME"))
	return GenerateJWT(&user, claims["auth_provider"].(string), accessExpiry, []byte(SecretKey))
}

// Logout function to revoke refresh token
func Logout(userID string) error {
	result := initializers.DB.Where("user_id = ?", userID).Delete(&models.RefreshToken{})
	return result.Error
}

const UserCtxKey = "user"

// Function to extract user role from context
func GetUserRoleFromJWT(ctx context.Context) (string, error) {
	claims, ok := ctx.Value(UserCtxKey).(jwt.MapClaims)
	if !ok {
		return "", errors.New("unauthorized")
	}
	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("role not found in token")
	}
	return role, nil
}

// Function to extract user from context
func GetUserFromJWT(ctx context.Context) (jwt.MapClaims, bool) {
	claims, ok := ctx.Value(UserCtxKey).(jwt.MapClaims)
	if !ok {
		fmt.Println("User not found in context")
		return nil, ok
	}
	name, ok := claims["name"].(string)
	if !ok {
		fmt.Println("Name not found in token")
	}
	fmt.Println("Name:", name)
	return claims, ok
}
