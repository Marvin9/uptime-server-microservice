package database

import (
	"fmt"
	"os"

	"github.com/Marvin9/uptime-server-microservice/pkg/models"

	"github.com/jinzhu/gorm"
)

// ConnectDB is used to connect database
func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", fmt.Sprintf("user=%v password=%v dbname=uptime_server_service sslmode=disable", os.Getenv("PSQL_USER"), os.Getenv("PSQL_PASSWORD")))
	if err != nil {
		return nil, err
	}

	database := db.DB()

	err = database.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

var migratedDbOnce = false

// InitialMigration to migrate initially
func InitialMigration(db *gorm.DB) {
	if !migratedDbOnce {
		db.AutoMigrate(&models.Users{}, &models.Reports{})
		migratedDbOnce = true
	}
}
