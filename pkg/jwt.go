package pkg

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Load secrets from environment variables
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))

func GenerateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Minute * 15).Unix(), // JWT expires in 15 minutes
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(), // Refresh token expires in 7 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

func ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	email, _ := claims["email"].(string)
	return email, nil
}
