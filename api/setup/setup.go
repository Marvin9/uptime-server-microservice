package setup

import (
	"net/http"

	"github.com/Marvin9/uptime-server-microservice/api"
	"github.com/Marvin9/uptime-server-microservice/api/middlewares"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/gin-gonic/gin"
)

// Router will setup APIs
func Router() *gin.Engine {
	r := gin.Default()

	authenticationGroup := r.Group("/auth")
	{
		authenticationGroup.POST("/register", api.RegisterAPI)
		authenticationGroup.POST("/login", api.LoginAPI)
	}

	authorizedGroup := r.Group("/api")
	{
		authorizedGroup.Use(middlewares.IsAuthorized())
		{
			authorizedGroup.GET("/", func(c *gin.Context) {
				jwtClaim, err := middlewares.ExtractJWTClaimFromContext(c)
				if err != nil {
					c.JSON(http.StatusUnauthorized, models.ErrorResponse(err.Error()))
					return
				}
				c.JSON(http.StatusOK, models.SuccessResponse(jwtClaim.UniqueID))
			})

			authorizedGroup.GET("/report", api.GetReportAPI)
			authorizedGroup.GET("/instances", api.GetInstancesAPI)
			authorizedGroup.POST("/instance", api.AddInstanceAPI)
			authorizedGroup.DELETE("/instance", api.RemoveInstanceAPI)
		}
	}

	return r
}
