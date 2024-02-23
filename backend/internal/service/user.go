package service

import (
	"strings"

	"day-trading-app/backend/internal/service/token"

	"golang.org/x/crypto/bcrypt"
)

func (s serviceImpl) AuthenticateUser(userName, password string) (string, error) {
	// 1. trim inputs
	userName = strings.TrimSpace(userName)

	// 2. get user from db
	user, err := s.db.GetUserByUserName(userName)
	if err != nil {
		return "", err
	}

	// 3. compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", err
	}

	// 4. generate and return jwt token
	jwt := token.NewJWTManager()
	token, err := jwt.GenerateToken(user.ID)
	return token, err
}
