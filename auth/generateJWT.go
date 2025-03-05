package auth

import (
	"fmt"
	"time"

	"github.com/Zenithive/it-crm-backend/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(user *models.User, authProvider string, expiryHours int, key []byte) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expiryHours) * time.Minute).Unix()

	fmt.Println("Refresh Key (generate):", key)

	fmt.Println("username:", user.Name)
	fmt.Println("role:", user.Role)

	claims := jwt.MapClaims{
		"user_id":       user.ID.String(),
		"name":          user.Name,
		"role":          user.Role,
		"auth_provider": authProvider,
		"exp":           expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("Token Signing Method:", token.Method)
	return token.SignedString(key)

}
