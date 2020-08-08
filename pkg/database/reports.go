package database

import (
	"net/http"

	"github.com/Marvin9/uptime-server-microservice/pkg/utils"

	"github.com/Marvin9/uptime-server-microservice/pkg/models"
)

// GetReports returns status code, error
func GetReports(userUniqueID string) (int, []models.Reports, error) {
	db, err := ConnectDB()
	defer db.Close()
	var reports = []models.Reports{}
	if err != nil {
		return http.StatusInternalServerError, reports, err
	}

	db.Preload("Instances", "Owner = ?", userUniqueID).Find(&reports)

	// db.Preload("Instance").Joins("JOIN instances ON reports.instance_id = instances.unique_id").Where("instances.owner = ?", userUniqueID).Find(&reports)

	// db.Select("url, status, reported_at, unique_id").Where("Owner = ?", userUniqueID).Order("reported_at DESC").Find(&reports)
	return http.StatusOK, reports, nil
}

// CreateInstance will create instance in db and return unique id of it.
func CreateInstance(userUniqueID, url string) (string, error) {
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
	}

	db.Create(&newInstance)

	return instanceUniqueID, nil
}
