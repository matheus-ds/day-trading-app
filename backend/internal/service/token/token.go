package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// NewJWTManager creates a new instance of JWTManager
func NewJWTManager() *JWTManager {
	// Retrieve the secret key from environment variable
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		panic("SECRET_KEY environment variable is not set")
	}

	//Hardcoded token duration to 24 hours
	tokenDuration := time.Hour * 24

	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// GenerateToken generates a new JWT token
func (jm *JWTManager) GenerateToken(userID *primitive.ObjectID) (string, error) {
	// Define the token claims
	claims := jwt.MapClaims{
		"user_ID": userID,                                  // Include any additional claims you need
		"exp":     time.Now().Add(jm.tokenDuration).Unix(), // Token expiration time
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(jm.secretKey))
	if err != nil {
		return "", err // Error generating token
	}

	return tokenString, nil
}
