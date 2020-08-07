package api

import (
	"log"
	"net/http"

	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/gin-gonic/gin"
)

// RegisterAPI is used to register user => /register
func RegisterAPI(c *gin.Context) {
	// request_body.go
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		log.Print("Error binding json.\n", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request body."))
		return
	}

	// TODO: EMAIL VALIDATION

	statusCode, err := database.RegisterUser(user.Email, user.Password)
	if err != nil {
		log.Println("Error registering user in database.\n", err)
		if statusCode != http.StatusInternalServerError {
			c.JSON(statusCode, models.ErrorResponse(err.Error()))
		} else {
			c.JSON(statusCode, models.ErrorResponse("Error registering user."))
		}
		return
	}

	c.JSON(statusCode, models.SuccessResponse("Successfully registered user."))
}
