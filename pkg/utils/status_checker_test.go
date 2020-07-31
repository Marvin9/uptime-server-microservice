package utils

import (
	"strconv"
	"testing"
)

var statusCodes = []int{
	200,
	404,
	400,
	503,
}

func statusCodeTestError(t *testing.T, expected, got int) {
	t.Errorf("Expected %v status code, got %v", expected, got)
}

func TestStatus(t *testing.T) {
	for _, expectedStatusCode := range statusCodes {
		url := "http://httpstat.us/" + strconv.Itoa(expectedStatusCode)
		gotStatusCode, _ := GetStatus(url)
		if gotStatusCode != expectedStatusCode {
			statusCodeTestError(t, expectedStatusCode, gotStatusCode)
		}
	}
}
