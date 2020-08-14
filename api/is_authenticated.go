package api

import (
	"net/http"

	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/gin-gonic/gin"
)

// IsAuthenticatedAPI is used to check if request is authenticated or not
func IsAuthenticatedAPI(c *gin.Context) {
	// This function will pass through middleware & is called only when request is authenticated
	c.JSON(http.StatusOK, models.SuccessResponse("Authenticated."))
}
