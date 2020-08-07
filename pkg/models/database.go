package models

import (
	"github.com/jinzhu/gorm"
)

// Users is table to store registered users
type Users struct {
	gorm.Model
	UniqueID string
	Email    string
	Password string
}
