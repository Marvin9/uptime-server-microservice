package setup

import (
	"flag"
	"net/http"

	"github.com/Marvin9/uptime-server-microservice/api"
	"github.com/Marvin9/uptime-server-microservice/api/middlewares"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Router will setup APIs
func Router() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	if flag.Lookup("test.v") != nil {
		config.AllowAllOrigins = true
	}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	r.Use(static.Serve("/", static.LocalFile("./client/dist", true)))

	authenticationGroup := r.Group("/auth")
	{
		authenticationGroup.POST("/register", api.RegisterAPI)
		authenticationGroup.POST("/login", api.LoginAPI)
		authenticationGroup.POST("/logout", api.LogoutAPI)

		authenticationGroup.Use(middlewares.IsAuthorized())
		{
			authenticationGroup.GET("/ping", api.IsAuthenticatedAPI)
		}
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
