package api

import (
	"net/http"
	"time"

	"github.com/Marvin9/uptime-server-microservice/api/middlewares"
	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/Marvin9/uptime-server-microservice/pkg/scheduler"
	"github.com/gin-gonic/gin"
)

const instanceDurationLowerBound = time.Hour

// AddInstanceAPI is used to add new server monitor => /api/instance
func AddInstanceAPI(c *gin.Context) {
	var instance models.Instance
	ok := instance.BindBody(c)
	if !ok {
		return
	}

	if instance.Duration < instanceDurationLowerBound {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Duration must be greater than or equal to 1 hour."))
		return
	}

	jwtClaims, err := middlewares.ExtractJWTClaimFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
		return
	}

	userUniqueID := jwtClaims.UniqueID
	if database.IsInstanceRunning(userUniqueID, instance.URL) {
		c.JSON(http.StatusConflict, models.ErrorResponse("There is already one instance for this URL."))
		return
	}

	newInstanceID, err := database.CreateInstance(userUniqueID, instance.URL, instance.Duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to create new instance"))
		return
	}

	statusCode, userEmail, err := database.GetUserEmail(userUniqueID)
	if err != nil {
		c.JSON(statusCode, models.ErrorResponse("Error creating instance. Try again."))
		return
	}

	go scheduler.InjectScheduler(newInstanceID, userUniqueID, userEmail, instance.URL, instance.Duration)

	c.JSON(http.StatusOK, models.SuccessResponse("Successfully added instance."))
}
