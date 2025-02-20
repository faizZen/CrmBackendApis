package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
		"user_id": user.ID.String(),
		"name":    user.Name,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token valid for 72 hours
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
