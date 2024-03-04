package middleware

import (
	"day-trading-app/backend/internal/service"
	"day-trading-app/backend/internal/service/token"
	"day-trading-app/backend/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(db service.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var userName string

		accessToken := c.Request.Header.Get("token")
		fmt.Println(accessToken)
		if userName, err = token.VerifyToken(accessToken); err == nil {
			fmt.Println(userName)
			// setting the authenticated userID, orgID in the context
			// we will use this info for logging and authorization purposes
			_, err := db.GetUserByUserName(userName)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
				logger.Error("Unauthorized request", logger.ErrorType(err))
				c.Abort()
				return
			}
			c.Set("user_name", userName)

			logger.Info("Authenticated user", logger.String("user_name", userName))

			c.Next()
			return

		}

		fmt.Println(err)
		logger.Info("Unauthorized request", logger.ErrorType(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
	}
}
