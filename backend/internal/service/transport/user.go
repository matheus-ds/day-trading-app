package transport

import (
	"net/http"

	"day-trading-app/backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

const AccessTokenDuration = 60 * 60 * 24 * 7 // 7 days

type AuthenticateUserReq struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func (e HTTPEndpoint) AuthenticateUser(c *gin.Context) {
	var user RegisterReq

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	token, err := e.srv.AuthenticateUser(user.UserName, user.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"data": gin.H{
				"error": err.Error(),
			},
		})
		logger.Error("failed to authenticate user", logger.ErrorType(err), logger.String("userName", user.UserName))

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"token": token,
		},
	})
	logger.Info("authenticated user", logger.String("userName", user.UserName))
}

type RegisterReq struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (e HTTPEndpoint) Register(c *gin.Context) {
	var user RegisterReq

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	err := e.srv.RegisterUser(user.UserName, user.Password, user.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    nil,
	})
}
