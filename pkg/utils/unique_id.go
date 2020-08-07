package utils

import (
	uuid "github.com/satori/go.uuid"
)

// GenerateUniqueID returns unique string to store as unique field in database
func GenerateUniqueID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
