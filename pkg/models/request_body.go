package models

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// User is for /auth/register post body
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// BindBody will map post body with User, if not same then it will response Bad request and return false
func (user *User) BindBody(c *gin.Context) bool {
	err := c.BindJSON(&user)
	if err != nil {
		log.Print("Error binding json.\n", err)
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body."))
		return false
	}
	return true
}

// Instance is for /api/instance post body
type Instance struct {
	URL      string        `json:"url"`
	Duration time.Duration `json:"duration"`
}

// BindBody will map post body with Instance
func (instance *Instance) BindBody(c *gin.Context) bool {
	err := c.BindJSON(&instance)
	if err != nil {
		log.Print("Error binding json.\n", err)
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body."))
		return false
	}
	return true
}
