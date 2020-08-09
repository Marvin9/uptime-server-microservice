package api

import (
	"log"
	"net/http"

	"github.com/Marvin9/uptime-server-microservice/api/middlewares"
	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/gin-gonic/gin"
)

// GetReportAPI is used to get the server uptimes of particular users => /api/report
func GetReportAPI(c *gin.Context) {
	// extract unique id to find all data of that user from table
	jwtClaim, err := middlewares.ExtractJWTClaimFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized"))
	}

	statusCode, reports, err := database.GetReports(jwtClaim.UniqueID)
	if err != nil {
		log.Println(err)
		c.JSON(statusCode, models.ErrorResponse("Error fetching report, Please try again."))
		return
	}

	response, err := models.SuccessResponseWithData(reports)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Error encoding json, Please try again."))
		return
	}

	c.JSON(http.StatusOK, response)
}
