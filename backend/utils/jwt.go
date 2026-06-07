package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key used to sign the JWT tokens. In production, this should be in an environment variable.
var jwtSecretKey = []byte("VaporAuror_Super_Secret_Key_2026")

// GenerateToken creates a JWT token for a given user ID and role.
func GenerateToken(userID uint, role string) (string, error) {
	// Create the JWT claims, which includes the user ID and expiry time
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token using the secret signed string
	return token.SignedString(jwtSecretKey)
}

// ParseToken validates the token string and extracts the claims.
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
