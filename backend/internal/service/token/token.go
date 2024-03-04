package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

type JWTManager struct {
	secretKey     []byte
	tokenDuration int
}

type AccessTokenClaims struct {
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

// NewJWTManager creates a new instance of JWTManager
func NewJWTManager() *JWTManager {
	// Retrieve the secret key from environment variable
	if os.Getenv("SECRET_KEY") == "" {
		panic("SECRET_KEY environment variable is not set")
	}

	//Hardcoded token duration to 24 hours
	tokenDuration := 24 * 60 * 60

	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// GenerateToken generates a new JWT token
func (jm *JWTManager) GenerateToken(userName string) (string, error) {
	// Define the token claims
	claims := AccessTokenClaims{
		userName, // Include any additional claims you need
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jm.tokenDuration) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "day-trading-app.com",
		},
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key and get the complete encoded token as a string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err // Error generating token
	}

	return tokenString, nil
}

func VerifyToken(accessToken string) (string, error) {
	var atc AccessTokenClaims

	token, err := jwt.ParseWithClaims(accessToken, &atc, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims.UserName, nil
	} else {
		return "", err
	}

}
