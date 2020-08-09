package database

import (
	"log"
	"net/http"
	"time"

	"github.com/Marvin9/uptime-server-microservice/pkg/utils"

	"github.com/Marvin9/uptime-server-microservice/pkg/models"
)

type schedulerStorage = map[string](map[string](chan bool))

// SchedulerList is all schedulers working
var SchedulerList = make(schedulerStorage)

// IsInstanceRunning returns true if that instance is already in memory
func IsInstanceRunning(ownerID, url string) bool {
	_, ok := SchedulerList[ownerID][url]
	return ok
}

// GetReports returns status code, error
func GetReports(userUniqueID string) (int, []models.ReportResponse, error) {
	db, err := ConnectDB()
	defer db.Close()
	var reports = []models.ReportResponse{}
	if err != nil {
		return http.StatusInternalServerError, reports, err
	}

	db.Raw("select instances.url, reports.* from reports, instances where reports.instance_id = instances.unique_id and instances.owner = ?", userUniqueID).Scan(&reports)

	return http.StatusOK, reports, nil
}

// CreateInstance will create instance in db and return unique id of it.
func CreateInstance(userUniqueID, url string, duration time.Duration) (string, error) {
	db, err := ConnectDB()
	defer db.Close()
	if err != nil {
		return "", err
	}

	instanceUniqueID, err := utils.GenerateUniqueID()
	if err != nil {
		return "", err
	}
	var newInstance = models.Instances{
		UniqueID: instanceUniqueID,
		Owner:    userUniqueID,
		URL:      url,
		Duration: duration,
	}

	db.Create(&newInstance)

	return instanceUniqueID, nil
}

// RemoveInstance will remove running instance from database and memory
func RemoveInstance(uniqueInstanceID string) error {
	db, err := ConnectDB()
	defer db.Close()
	if err != nil {
		return err
	}

	var instance = models.Instances{
		UniqueID: uniqueInstanceID,
	}

	db.Find(&instance)
	if IsInstanceRunning(instance.Owner, instance.URL) {
		SchedulerList[instance.Owner][instance.URL] <- true
		delete(SchedulerList[instance.Owner], instance.URL)
	}

	err = db.Exec("DELETE FROM instances WHERE unique_id = ?;", uniqueInstanceID).Error
	if err != nil {
		log.Println(err)
	}

	err = db.Exec("DELETE FROM reports WHERE instance_id = ?;", uniqueInstanceID).Error
	if err != nil {
		log.Println(err)
	}

	return nil
}
