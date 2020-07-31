package scheduler

import (
	"log"
	"sync"
	"time"

	"github.com/Marvin9/uptime-server-microservice/pkg/utils"
)

// Schedule - schedule to check status of url after delay
func Schedule(forurl string, delay time.Duration, wg *sync.WaitGroup, action func(string, time.Time, int)) chan bool {
	wg.Add(1)
	defer wg.Done()
	stop := make(chan bool)

	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			status, err := utils.GetStatus(forurl)
			if err != nil {
				log.Fatal(err)
				return
			}
			action(forurl, time.Now(), status)
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}
