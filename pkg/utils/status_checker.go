package utils

import (
	"log"
	"net/http"
)

// GetStatus - gives the status of passed url.
func GetStatus(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("%v\n", err)
		return -1, err
	}
	return resp.StatusCode, nil
}
