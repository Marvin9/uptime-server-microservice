package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Marvin9/uptime-server-microservice/pkg/scheduler"
)

const url string = "https://marvin9-web-analyzer-server.glitch.me/"

var urls = []string{
	"https://marvin9-web-analyzer-server.glitch.me",
	"https://www.google.com",
	"https://www.facebook.com",
}

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
	var schedulers = make(schedulerStorage)
	for _, url := range urls {
		injectScheduler(url, 5*time.Second, schedulers)
	}
	waitGroup.Wait()
}
