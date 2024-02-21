package service

import (
	"strings"
	"time"

	token "github.com/matheus-ds/day-trading-app/backend/internal/service/jwt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64  `db:"id"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
}

func (s serviceImpl) AuthenticateUser(email, password string) (string, error) {
	// TODO: implement this

	// 1. trim inputs
	email = strings.TrimSpace(email)

	// 2. get user from db
	user := User{}
	// user, err := s.db.GetUserByEmail(email)

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", err
	}

	// 3. generate and return jwt token
	jwt := token.NewJWTManager(time.Hour)
	token, err := jwt.GenerateToken(user.ID)
	return token, err
}
