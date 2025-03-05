package auth

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Zenithive/it-crm-backend/models"
)

func GenerateTokens(user *models.User, authProvider string) (string, string, error) {
	// Access Token (Short-lived)
	accessExpiry, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_TIME")) // e.g., 15 min
	accessToken, err := GenerateJWT(user, authProvider, accessExpiry, []byte(SecretKey))
	if err != nil {
		return "", "", err
	}

	// Refresh Token (Long-lived)
	refreshExpiry, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY")) // e.g., 7 days
	refreshToken, err := GenerateJWT(user, authProvider, refreshExpiry, []byte(RefreshSecretKey))
	if err != nil {
		return "", "", err
	}
	fmt.Println("Access Token:", accessToken)
	fmt.Println("Refresh Token:", refreshToken)
	fmt.Println("Access Expiry:", accessExpiry)
	fmt.Println("Refresh Expiry:", refreshExpiry)

	// Store refresh token in the database
	err = StoreRefreshToken(user.ID.String(), refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
