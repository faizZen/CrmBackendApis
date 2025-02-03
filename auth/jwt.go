package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Zenithive/it-crm-backend/models"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("abcdefghijklmnopqrstuvwxyz") // Change this to a secure key

// Claims structure for JWT
type Claims struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a new token
func GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"name":    user.Name,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateJWT validates the token and extracts claims
func ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	fmt.Println("Validating Token:", tokenStr) // Debug log

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		fmt.Println("JWT Parsing Error:", err) // Debug log
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		fmt.Println("Invalid JWT Claims")
		return nil, errors.New("invalid token")
	}

	fmt.Println("JWT is valid. Extracted Claims:", claims)
	return claims, nil
}

const UserCtxKey = "user"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		fmt.Println("Authorization Header:", authHeader) // Debug log

		if authHeader == "" {
			fmt.Println("No Authorization header found")
			next.ServeHTTP(w, r)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println("Extracted Token:", tokenString) // Debug log

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			fmt.Println("JWT Validation Error:", err) // Debug log
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		fmt.Println("Token is valid. Claims:", claims)
		ctx := context.WithValue(r.Context(), UserCtxKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext extracts the user claims from the request context
func GetUserFromContext(ctx context.Context) (jwt.MapClaims, bool) {
	claims, ok := ctx.Value(UserCtxKey).(jwt.MapClaims)
	name, ok := claims["name"].(string)
	if !ok {
		// http.Error(, "Name not found in token", http.StatusUnauthorized)
	}
	fmt.Println("Name:", name)
	return claims, ok
}
