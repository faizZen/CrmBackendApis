package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Zenithive/it-crm-backend/models"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("SECRET_KEY")) // Change this to a secure key

// Claims structure for JWT
type Claims struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
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
		return secretKey, nil
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
func GenerateJWT(user *models.User) (string, error) {
	expiryHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_TIME")) // Convert string to int
	expirationTime := time.Now().Add(time.Duration(expiryHours) * time.Hour).Unix()
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"name":    user.Name,
		"role":    user.Role,
		// "auth_provider": authProvider, // "local" or "google"
		"exp": expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateJWT validates the token and extracts claims
func ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims in token")
	}

	return claims, nil
}

const UserCtxKey = "user"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Read request body to extract GraphQL operation name
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		/*
			!Important
			When you call:
			body, err := ioutil.ReadAll(r.Body)
				It reads the entire request body.
				But ReadAll consumes the body, meaning it cannot be read again later in the request lifecycle.
				If you try to access r.Body later (e.g., in another middleware or handler), it will be empty.
				strings.NewReader(string(body)) → Creates a new readable stream from the body.
				io.NopCloser(...) → Wraps it in a no-op closer, so r.Body.Close() doesn't break anything.
		*/

		r.Body = io.NopCloser(strings.NewReader(string(body)))
		// Parse JSON request body
		var graphqlReq struct {
			Query string `json:"query"`
		}
		err = json.Unmarshal(body, &graphqlReq)
		if err != nil {
			http.Error(w, "Invalid GraphQL request format", http.StatusBadRequest)
			return
		}

		// *Allow login and register mutations without a token

		if strings.Contains(graphqlReq.Query, "login") {
			fmt.Println("Login mutation detected, skipping auth check.")
			next.ServeHTTP(w, r)
			return
		}

		//* Check token for all Other mutations

		authHeader := r.Header.Get("Authorization")
		fmt.Println("Authorization Header:", authHeader)

		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Store claims in request context and pass to next handler
		ctx := context.WithValue(r.Context(), UserCtxKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func MiddlewareFuncForUploads(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(token, "Bearer ")
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserCtxKey, claims)
		next(w, r.WithContext(ctx))
	}
}

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
