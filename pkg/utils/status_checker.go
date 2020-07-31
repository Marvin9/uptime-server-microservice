package utils

import (
	"log"
	"net/http"
)

// GetStatus - gives the status of passed url.
func GetStatus(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return -1, err
	}
	return resp.StatusCode, nil
}
