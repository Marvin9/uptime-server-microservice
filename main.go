package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
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
	// var schedulers = make(schedulerStorage)
	r := gin.Default()
	r.POST("/register", func(c *gin.Context) {
		// request_body.go
		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			log.Print("Error binding json.\n", err)
			c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request body."))
			return
		}

		// TODO: EMAIL VALIDATION

		statusCode, err := database.RegisterUser(user.Email, user.Password)
		if err != nil {
			log.Println("Error registering user in database.\n", err)
			if statusCode != http.StatusInternalServerError {
				c.JSON(statusCode, models.ErrorResponse(err.Error()))
			} else {
				c.JSON(statusCode, models.ErrorResponse("Error registering user."))
			}
			return
		}

		c.JSON(statusCode, models.SuccessResponse("Successfully registered user."))
	})
	r.Run(":8000")
}
