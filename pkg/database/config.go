package database

import (
	"log"
	"os"

	"github.com/Marvin9/uptime-server-microservice/pkg/models"

	"github.com/jinzhu/gorm"
)

// ConnectDB is used to connect database
func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Printf("Inside ConnectDB: %v\n\n", err)
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
		db.AutoMigrate(&models.Users{}, &models.Reports{}, &models.Instances{})
		migratedDbOnce = true
	}
}
