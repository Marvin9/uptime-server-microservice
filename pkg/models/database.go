package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Users is table to store registered users
type Users struct {
	gorm.Model
	UniqueID string `gorm:"primary_key"`
	Email    string
	Password string
}

// Reports is table to store report of server status
type Reports struct {
	gorm.Model
	UniqueID   string `gorm:"primary_key"`
	Owner      string
	URL        string
	Status     int
	ReportedAt time.Time
}
