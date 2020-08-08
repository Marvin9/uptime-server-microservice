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

type schedulerStorage = map[string](map[string](chan bool))

// SchedulerList is all schedulers working
var SchedulerList = make(schedulerStorage)

// Action will be envoked after certain duration
func Action(newInstanceID, url string, t time.Time, status int) {
	db, err := database.ConnectDB()
	if err != nil {
		log.Print("Error connecting database.\n", err)
		return
	}
	defer db.Close()

	// GET LATEST REPORT
	var latestReport = models.Reports{}
	db.Where("instance_id = ?", newInstanceID).Order("reported_at DESC").First(&latestReport)

	// IF THE STATUS FLUCTUATE THEN STORE IT IN DATABASE
	if latestReport.Status != status {
		uniqueID, err := utils.GenerateUniqueID()
		if err != nil {
			log.Print("Error generating unique id.\n", err)
			return
		}

		instance := models.Reports{
			UniqueID:   uniqueID,
			InstanceID: newInstanceID,
			Status:     status,
			ReportedAt: t,
		}

		db.Create(&instance)
	}
}

// InjectScheduler is used to add instance for server monitoring
func InjectScheduler(newInstanceID, ownerUniqueID, url string, delay time.Duration) bool {
	waitGroup.Add(1)
	_, ok := SchedulerList[ownerUniqueID][url]
	if !ok {
		stop := Schedule(newInstanceID, url, delay)
		_, ook := SchedulerList[ownerUniqueID]
		if !ook {
			SchedulerList[ownerUniqueID] = make(map[string](chan bool))
		}
		SchedulerList[ownerUniqueID][url] = stop
	}
	return !ok
}

// Schedule - schedule to check status of url after delay
func Schedule(newInstanceID, forurl string, delay time.Duration) chan bool {
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
			Action(newInstanceID, forurl, time.Now(), status)
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

// IsInstanceRunning returns true if that instance is already in memory
func IsInstanceRunning(ownerID, url string) bool {
	_, ok := SchedulerList[ownerID][url]
	return ok
}
