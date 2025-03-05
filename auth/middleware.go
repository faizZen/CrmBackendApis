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

	initializers "github.com/Zenithive/it-crm-backend/Initializers"
	"github.com/Zenithive/it-crm-backend/models"
	"github.com/golang-jwt/jwt/v4"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Read request body to extract GraphQL operation name
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(strings.NewReader(string(body)))

		var graphqlReq struct {
			Query string `json:"query"`
		}
		err = json.Unmarshal(body, &graphqlReq)
		if err != nil {
			http.Error(w, "Invalid GraphQL request format", http.StatusBadRequest)
			return
		}

		// Allow login mutation without a token
		if strings.Contains(graphqlReq.Query, "login") {
			fmt.Println("Login mutation detected, skipping auth check.")
			next.ServeHTTP(w, r)
			return
		}

		// Check token for all other requests
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateJWT(tokenString, []byte(SecretKey))

		if err != nil {
			fmt.Println("error validating token:", err)

			var validationErr *jwt.ValidationError
			if errors.As(err, &validationErr) {
				fmt.Printf("ValidationError detected: %v\n", validationErr.Errors)

				if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
					fmt.Println("Access token expired. Attempting refresh...")

					// Parse token without validation to extract claims
					parsedToken, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
					if err != nil {
						http.Error(w, "Unauthorized: Unable to parse token", http.StatusUnauthorized)
						return

					}

					// Extract user ID
					claims, ok := parsedToken.Claims.(jwt.MapClaims)
					if !ok {
						http.Error(w, "Unauthorized: Invalid token claims", http.StatusUnauthorized)
						return
					}

					userID, ok := claims["user_id"].(string)
					if !ok {
						http.Error(w, "Unauthorized: Invalid user ID in token", http.StatusUnauthorized)
						return
					}

					// Fetch refresh token from DB
					var refreshRecord models.RefreshToken
					result := initializers.DB.Where("user_id = ?", userID).First(&refreshRecord)
					if result.Error != nil {
						http.Error(w, "Unauthorized: No valid refresh token found", http.StatusUnauthorized)
						return
					}
					fmt.Println("Refresh token found:", refreshRecord.Token)

					// Validate refresh token
					newClaims, err := ValidateJWT(refreshRecord.Token, []byte(RefreshSecretKey))
					if err != nil {
						fmt.Println("Refresh Token Validation Error:", err) // Log error
						http.Error(w, "Unauthorized: Invalid refresh token", http.StatusUnauthorized)
						return
					}

					// Retrieve user info
					var user models.User
					result = initializers.DB.Where("id = ?", userID).First(&user)
					if result.Error != nil {
						http.Error(w, "Unauthorized: User not found", http.StatusUnauthorized)
						return
					}
					fmt.Println("User found:", user.Name)
					fmt.Println("User Role:", user.Role)
					// Generate a new access token
					accessExpiry, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_TIME"))
					newAccessToken, err := GenerateJWT(&user, newClaims["auth_provider"].(string), accessExpiry, []byte(SecretKey))
					if err != nil {
						http.Error(w, "Failed to generate new access token", http.StatusInternalServerError)
						return
					}

					// Attach the new access token in the response header
					w.Header().Set("New-Access-Token", newAccessToken)
					fmt.Println("New access token issued:", newAccessToken)

					fmt.Println("Final Claims Before Storing in Context with expired access:", newClaims)
					// Continue request processing with updated claims
					ctx := context.WithValue(r.Context(), UserCtxKey, newClaims)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			} else if strings.Contains(err.Error(), "Token is expired") {
				// Fallback: Check error message manually if wrapped
				fmt.Println("Access token expired (detected via string match). Attempting refresh...")

				// Parse token without validation to extract claims
				parsedToken, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
				if err != nil {
					http.Error(w, "Unauthorized: Unable to parse token", http.StatusUnauthorized)
					return

				}

				// Extract user ID
				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				if !ok {
					http.Error(w, "Unauthorized: Invalid token claims", http.StatusUnauthorized)
					return
				}

				userID, ok := claims["user_id"].(string)
				if !ok {
					http.Error(w, "Unauthorized: Invalid user ID in token", http.StatusUnauthorized)
					return
				}

				// Fetch refresh token from DB
				var refreshRecord models.RefreshToken
				result := initializers.DB.Where("user_id = ?", userID).First(&refreshRecord)
				if result.Error != nil {
					http.Error(w, "Unauthorized: No valid refresh token found", http.StatusUnauthorized)
					return
				}
				fmt.Println("Refresh token found:", refreshRecord.Token)

				// Validate refresh token
				newClaims, err := ValidateRefreshToken(refreshRecord.Token)
				if err != nil {
					fmt.Println("Refresh Token Validation Error:", err) // Log error
					http.Error(w, "Unauthorized: Invalid refresh token", http.StatusUnauthorized)
					return
				}

				// Retrieve user info
				var user models.User
				result = initializers.DB.Where("id = ?", userID).First(&user)
				if result.Error != nil {
					http.Error(w, "Unauthorized: User not found", http.StatusUnauthorized)
					return
				}
				fmt.Println("User found:", user.Name)
				fmt.Println("User Role:", user.Role)
				// Generate a new access token
				accessExpiry, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_TIME"))
				newAccessToken, err := GenerateJWT(&user, newClaims["auth_provider"].(string), accessExpiry, []byte(SecretKey))
				if err != nil {
					http.Error(w, "Failed to generate new access token", http.StatusInternalServerError)
					return
				}

				// Attach the new access token in the response header
				w.Header().Set("New-Access-Token", newAccessToken)
				fmt.Println("New access token issued:", newAccessToken)

				fmt.Println("Final Claims Before Storing in Context with expired access:", newClaims)
				// Continue request processing with updated claims
				ctx := context.WithValue(r.Context(), UserCtxKey, newClaims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			} else {
				fmt.Println("Unexpected error:", err)
				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}

		}
		fmt.Println("Final Claims Before Storing in Context with accesstoken:", claims)

		// If token is valid, continue processing the request
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
		claims, err := ValidateJWT(tokenString, []byte(SecretKey))
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserCtxKey, claims)
		next(w, r.WithContext(ctx))
	}
}
