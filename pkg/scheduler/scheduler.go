package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Marvin9/uptime-server-microservice/pkg/mailer"

	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"

	"github.com/Marvin9/uptime-server-microservice/pkg/utils"
)

var waitGroup sync.WaitGroup

// Action will be envoked after certain duration
func Action(newInstanceID, userEmail, url string, t time.Time, status int) {
	fmt.Println("\nAction for instance: ", newInstanceID, ", URL: ", url)
	fmt.Println()
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
		go mailer.Mail(userEmail, url, status, t)
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
func InjectScheduler(newInstanceID, ownerUniqueID, userEmail, url string, delay time.Duration) bool {
	waitGroup.Add(1)
	_, ok := database.SchedulerList[ownerUniqueID][url]
	if !ok {
		stop := Schedule(newInstanceID, userEmail, url, delay)
		_, ook := database.SchedulerList[ownerUniqueID]
		if !ook {
			database.SchedulerList[ownerUniqueID] = make(map[string](chan bool))
		}
		database.SchedulerList[ownerUniqueID][url] = stop
	}
	return !ok
}

// Schedule - schedule to check status of url after delay
func Schedule(newInstanceID, userEmail, forurl string, delay time.Duration) chan bool {
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
			Action(newInstanceID, userEmail, forurl, time.Now(), status)
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}
