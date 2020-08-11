package scheduler

import (
	"log"

	"github.com/Marvin9/uptime-server-microservice/pkg/database"
	"github.com/Marvin9/uptime-server-microservice/pkg/models"
)

// RestartAllSchedulers is used to restart scheduler whenever server starts
func RestartAllSchedulers() {
	if len(database.SchedulerList) == 0 {
		var instances []models.Instances
		db, err := database.ConnectDB()
		defer db.Close()
		if err != nil {
			log.Println("Error connecting database.\n", err)
			return
		}
		db.Find(&instances)

		for _, instance := range instances {
			_, userEmail, err := database.GetUserEmail(instance.Owner)
			if err != nil {
				log.Printf("Skipped scheduler for instance %v.", instance.Owner)
				continue
			}
			go InjectScheduler(instance.UniqueID, instance.Owner, userEmail, instance.URL, instance.Duration)
		}

		log.Println("Restarted all schedulers")
	}
}
