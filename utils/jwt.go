package utils

import (
	"location-share-backend/initializers"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(sessionId string) (string, error) {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"session_id": sessionId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
