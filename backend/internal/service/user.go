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
	token, err := jwt.GenerateToken(user.UserName)
	return token, err
}

func (s serviceImpl) RegisterUser(userName, password, name string) error {
	userName = strings.TrimSpace(userName)

	_, err := s.db.GetUserByUserName(userName)
	if err == nil {
		return err // user already exists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.db.RegisterUser(userName, string(hashedPassword))
	return err
}
