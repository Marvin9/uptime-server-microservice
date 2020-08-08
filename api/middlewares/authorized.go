package middlewares

import (
	"net/http"

	"github.com/Marvin9/uptime-server-microservice/pkg/utils"
	"github.com/dgrijalva/jwt-go"

	"github.com/Marvin9/uptime-server-microservice/api"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/gin-gonic/gin"
)

const unauthorized = "Unauthorized"

// IsAuthorized is middleware to check if user is logged in or not
func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtCookie, err := c.Cookie(api.JWTCookieName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse(unauthorized))
			c.Abort()
			return
		}

		tokenString := jwtCookie

		jwtClaims := &models.Claims{}

		signedToken, err := jwt.ParseWithClaims(tokenString, jwtClaims, func(token *jwt.Token) (interface{}, error) {
			return utils.GetJWTKey(), nil
		})

		if err != nil || !signedToken.Valid {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse(unauthorized))
			c.Abort()
			return
		}

		c.Next()
	}
}
