package api

import (
	"net/http"

	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/gin-gonic/gin"
)

// RemoveInstanceAPI will remove instance created by user
func RemoveInstanceAPI(c *gin.Context) {
	var uniqueInstance models.UniqueInstance
	ok := uniqueInstance.BindBody(c)
	if !ok {
		return
	}

	err := database.RemoveInstance(uniqueInstance.InstanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Could not remove instance. Please try again."))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Successfully removed instance."))
}
