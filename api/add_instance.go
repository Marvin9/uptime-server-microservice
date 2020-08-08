package api

import (
	"net/http"
	"time"

	"github.com/Marvin9/uptime-server-microservice/api/middlewares"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/Marvin9/uptime-server-microservice/pkg/scheduler"
	"github.com/gin-gonic/gin"
)

const instanceDurationLowerBound = time.Hour

// AddInstance is used to add new server monitor => /api/instance
func AddInstance(c *gin.Context) {
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
	schedulerAdded := scheduler.InjectScheduler(userUniqueID, instance.URL, instance.Duration)

	if schedulerAdded {
		c.JSON(http.StatusOK, models.SuccessResponse("Successfully added instance."))
	} else {
		c.JSON(http.StatusConflict, models.ErrorResponse("There is already one instance for this URL."))
	}
}
