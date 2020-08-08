package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Marvin9/uptime-server-microservice/pkg/database"

	"github.com/Marvin9/uptime-server-microservice/pkg/models"

	"github.com/Marvin9/uptime-server-microservice/api"
	"github.com/Marvin9/uptime-server-microservice/api/middlewares"
	"github.com/Marvin9/uptime-server-microservice/pkg/scheduler"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

var waitGroup sync.WaitGroup

type schedulerStorage = map[string](chan bool)

func action(url string, t time.Time, status int) {
	fmt.Printf("Status of %v at %v: %v\n\n", url, t, status)
}

func injectScheduler(url string, delay time.Duration, schedulerList schedulerStorage) {
	waitGroup.Add(1)
	stop := scheduler.Schedule(url, delay, &waitGroup, action)
	schedulerList[url] = stop
}

func main() {
	db, _ := database.ConnectDB()
	database.InitialMigration(db)
	db.Close()

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

			authorizedGroup.GET("/report", api.GetReport)
		}
	}

	r.Run(":8000")
}
