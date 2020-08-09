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

func bindingError(c *gin.Context, err error) bool {
	if err != nil {
		log.Print("Error binding json.\n", err)
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body."))
		return false
	}
	return true
}

// BindBody will map post body with User, if not same then it will response Bad request and return false
func (user *User) BindBody(c *gin.Context) bool {
	err := c.BindJSON(&user)
	return bindingError(c, err)
}

// Instance is for /api/instance post body
type Instance struct {
	URL      string        `json:"url"`
	Duration time.Duration `json:"duration"`
}

// BindBody will map post body with Instance
func (instance *Instance) BindBody(c *gin.Context) bool {
	err := c.BindJSON(&instance)
	return bindingError(c, err)
}

// UniqueInstance is for /api/instance delete body
type UniqueInstance struct {
	InstanceID string `json:"instance_id"`
}

// BindBody will map delete body with UniqueInstance
func (instance *UniqueInstance) BindBody(c *gin.Context) bool {
	err := c.BindJSON(&instance)
	return bindingError(c, err)
}
