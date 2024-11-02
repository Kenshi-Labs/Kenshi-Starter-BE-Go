package utils

import (
	"time"

	"auth-api/configs"
	"auth-api/models"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a new JWT token for a user
func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"id": user.ID.Hex(),
		"email": user.Email,
		"role":user.Role,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 24 hours expiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configs.AppConfig.JWTSecret))
}

// ValidateJWT verifies the JWT token and returns the claims
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.AppConfig.JWTSecret),nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}