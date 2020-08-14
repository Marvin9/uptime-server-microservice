package api

import (
	"net/http"

	"github.com/Marvin9/uptime-server-microservice/pkg/database"

	"github.com/Marvin9/uptime-server-microservice/api/middlewares"

	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/gin-gonic/gin"
)

// IsAuthenticatedAPI is used to check if request is authenticated or not
func IsAuthenticatedAPI(c *gin.Context) {
	jwtClaims, err := middlewares.ExtractJWTClaimFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Error extracting jwt."))
		return
	}

	statusCode, mail, err := database.GetUserEmail(jwtClaims.UniqueID)
	if err != nil {
		c.JSON(statusCode, models.ErrorResponse("Error fetching registered user from database."))
		return
	}

	response, err := models.SuccessResponseWithData(models.PingResponse{
		Email: mail,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Error generating response."))
	}

	c.JSON(statusCode, response)
}
