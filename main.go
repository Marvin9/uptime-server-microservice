package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Marvin9/uptime-server-microservice/api"
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

	authenticationGroup := r.Group("/auth")
	{
		authenticationGroup.POST("/register", api.RegisterAPI)
		authenticationGroup.POST("/login", api.LoginAPI)
	}

	r.Run(":8000")
}
