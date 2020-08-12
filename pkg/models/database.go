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
	UniqueID   string `gorm:"primary_key" json:"unique_id"`
	InstanceID string `json:"instance_id"`
	Instace    Instances
	Status     int       `json:"status"`
	ReportedAt time.Time `json:"reported_at"`
}

// Instances is table to store instance with map of => owner & url
type Instances struct {
	gorm.Model
	UniqueID string        `gorm:"primary_key" json:"unique_id,omitempty"`
	Owner    string        `json:",omitempty"`
	URL      string        `json:"url"`
	Duration time.Duration `json:"duration"`
}
