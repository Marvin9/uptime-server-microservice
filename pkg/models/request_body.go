package models

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// User is for /register post body
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
