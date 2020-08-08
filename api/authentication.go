package api

import (
	"log"

	"github.com/Marvin9/uptime-server-microservice/pkg/utils"

	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/gin-gonic/gin"
)

const JWTCookieName = "jwt"

// RegisterAPI is used to register user => /auth/register
func RegisterAPI(c *gin.Context) {
	// request_body.go
	var user models.User
	ok := user.BindBody(c)
	if !ok {
		return
	}

	// TODO: EMAIL VALIDATION

	statusCode, err := database.RegisterUser(user.Email, user.Password)
	if err != nil {
		log.Println("Error registering user in database.\n", err)
		c.JSON(statusCode, models.CustomErrorResponse(statusCode, err, "Error registering user."))
		return
	}

	c.JSON(statusCode, models.SuccessResponse("Successfully registered user."))
}

// LoginAPI is used to login user => /auth/login
func LoginAPI(c *gin.Context) {
	var user models.User
	ok := user.BindBody(c)
	if !ok {
		return
	}

	statusCode, jwtToken, err := database.LoginUser(user.Email, user.Password)
	if err != nil {
		log.Println("Error logging in.\n", err)
		c.JSON(statusCode, models.CustomErrorResponse(statusCode, err, "Error logging in."))
		return
	}

	c.SetCookie(JWTCookieName, jwtToken, int(utils.JWTExpireAfter), "/", "localhost", false, true)
	c.JSON(statusCode, models.SuccessResponse("Successfully logged in."))
}
