package utils

import (
	uuid "github.com/satori/go.uuid"
)

// GenerateUniqueID returns unique string to store as unique field in database
func GenerateUniqueID() (string, error) {
	id := uuid.NewV4()

	return id.String(), nil
}
