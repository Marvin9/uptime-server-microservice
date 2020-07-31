package utils

import "log"

// LogError - error handling utility
func LogError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
