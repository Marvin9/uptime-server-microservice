package scheduler

import (
	"log"
	"sync"
	"time"

	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"

	"github.com/Marvin9/uptime-server-microservice/pkg/utils"
)

var waitGroup sync.WaitGroup

type schedulerStorage = map[string](chan bool)

// SchedulerList is all schedulers working
var SchedulerList = make(schedulerStorage)

// Action will be envoked after certain duration
func Action(ownerUniqueID, url string, t time.Time, status int) {
	db, err := database.ConnectDB()
	if err != nil {
		log.Print("Error connecting database.\n", err)
		return
	}
	defer db.Close()

	uniqueID, err := utils.GenerateUniqueID()
	if err != nil {
		log.Print("Error generating unique id.\n", err)
		return
	}

	instance := models.Reports{
		UniqueID:   uniqueID,
		Owner:      ownerUniqueID,
		URL:        url,
		Status:     status,
		ReportedAt: t,
	}

	db.Create(&instance)
}

// InjectScheduler is used to add instance for server monitoring
func InjectScheduler(ownerUniqueID, url string, delay time.Duration) {
	waitGroup.Add(1)
	_, ok := SchedulerList[url]
	if !ok {
		stop := Schedule(ownerUniqueID, url, delay)
		SchedulerList[url] = stop
	}
}

// Schedule - schedule to check status of url after delay
func Schedule(ownerUniqueID, forurl string, delay time.Duration) chan bool {
	waitGroup.Add(1)
	defer waitGroup.Done()
	stop := make(chan bool)

	go func() {
		waitGroup.Add(1)
		defer waitGroup.Done()
		for {
			status, err := utils.GetStatus(forurl)
			if err != nil {
				log.Print(err)
				log.Print("\n")
				return
			}
			Action(ownerUniqueID, forurl, time.Now(), status)
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}
