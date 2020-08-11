package main

import (
	"net/http"

	"github.com/Marvin9/uptime-server-microservice/api"
	"github.com/Marvin9/uptime-server-microservice/api/middlewares"
	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
	"github.com/Marvin9/uptime-server-microservice/pkg/scheduler"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	db, _ := database.ConnectDB()
	database.InitialMigration(db)
	db.Close()

	go scheduler.RestartAllSchedulers()

	// var schedulers = make(schedulerStorage)
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

	r.Run(":8000")
}
