package utils

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("your_secret_key") // Replace this in production!

// GenerateToken creates a JWT for a user
func GenerateToken(userID uint) (string, error) {
	exp := time.Now().Add(time.Hour * 24 * 7).Unix()
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     exp,
	}

	log.Printf("Generated token for user %d with expiration: %d\n", userID, exp)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		log.Printf("Token parse error: %v\n", err)
		return nil, err
	}

	if !token.Valid {
		log.Println("Token is invalid")
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Failed to parse claims")
		return nil, err
	}

	log.Printf("Parsed token claims: %v\n", claims)
	return claims, nil
}
