package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func ValidateJWT(tokenStr string, key []byte) (jwt.MapClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("missing token")
	}

	// Ensure token has three segments (header.payload.signature)
	if len(strings.Split(tokenStr, ".")) != 3 {
		return nil, errors.New("malformed token")
	}
	// Parse the token with error handling
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims in token")
	}

	// Validate expiration
	exp, ok := claims["exp"].(float64) // JWT stores exp as float64
	if !ok {
		return nil, errors.New("invalid exp format")
	}

	if int64(exp) < time.Now().Unix() {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
