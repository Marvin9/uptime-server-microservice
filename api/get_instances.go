package api

import (
	"net/http"

	"github.com/Marvin9/uptime-server-microservice/api/middlewares"
	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/gin-gonic/gin"
)

// GetInstancesAPI is used get current instances use has started => /api/instances
func GetInstancesAPI(c *gin.Context) {
	jwtClaims, err := middlewares.ExtractJWTClaimFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	userUniqueID := jwtClaims.UniqueID

	statusCode, instances, err := database.GetInstances(userUniqueID)

	if err != nil {
		c.JSON(statusCode, models.ErrorResponse("Error while fetching instances"))
		return
	}

	response, _ := models.SuccessResponseWithData(instances)

	c.JSON(statusCode, response)
}
