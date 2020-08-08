package database

import (
	"net/http"

	"github.com/Marvin9/uptime-server-microservice/pkg/models"
)

// GetReports returns status code, error
func GetReports(userUniqueID string) (int, []models.Reports, error) {
	db, err := ConnectDB()
	var reports = []models.Reports{}
	if err != nil {
		return http.StatusInternalServerError, reports, err
	}
	defer db.Close()

	db.Select("url, status, reported_at").Where("Owner = ?", userUniqueID).Order("reported_at DESC").Find(&reports)
	return http.StatusOK, reports, nil
}
