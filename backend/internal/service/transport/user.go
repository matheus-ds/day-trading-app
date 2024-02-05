package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheus-ds/day-trading-app/backend/pkg/logger"
)

const AccessTokenDuration = 60 * 60 * 24 * 7 // 7 days

func (e HTTPEndpoint) AuthenticateUser(c *gin.Context) {
	// TODO: grab from request payload
	email := "example"
	password := "example"

	token, err := e.srv.AuthenticateUser(email, password)

	if err == nil {
		c.SetCookie("access_token", string(token), AccessTokenDuration, "/", "", true, true)

		c.JSON(http.StatusOK, gin.H{})
		logger.Info("authenticated user")

		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "failed to authenticate user",
		"email": email,
	})
	logger.Error("failed to authenticate user", logger.ErrorType(err), logger.String("email", email))
}

func (e HTTPEndpoint) Register(c *gin.Context) {
	// TODO: implement

	c.JSON(http.StatusOK, gin.H{})
}
