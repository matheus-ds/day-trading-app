package middleware

import (
	"net/http"
	"strings"

	"github.com/matheus-ds/day-trading-app/authentication"
	"github.com/matheus-ds/day-trading-app/tokenutil" //not sure if i put in the othere 2 go files correctly? the package name at the top all the same? or different plz help me
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, err := tokenutil.IsAuthorized(authToken, secret)
			if authorized {
				userID, err := tokenutil.ExtractIDFromToken(authToken, secret)
				if err != nil {
					c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
					c.Abort()
					return
				}
				c.Set("x-user-id", userID)
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Not authorized"})
		c.Abort()
	}
}


//We can get the UserID from the HTTP Web Framework Context as below:
userID := c.GetString("x-user-id")

//Then, we can use this middleware as below:
router.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))

//When we want to generate the access token, we can call:
accessToken, err := tokenutil.CreateAccessToken(user, secret, expiry)

//For generating the refresh token, we can call:
refreshToken, err := tokenutil.CreateRefreshToken(user, secret, expiry)