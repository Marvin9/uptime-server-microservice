package models

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorResponse is response of any request, if any error
func ErrorResponse(err string) gin.H {
	return gin.H{
		"error":   true,
		"message": err,
	}
}

// CustomErrorResponse will return default error message if it is not internal error.
func CustomErrorResponse(responseStatusCode int, err error, customMessage string) gin.H {
	if responseStatusCode != http.StatusInternalServerError {
		return ErrorResponse(err.Error())
	}
	return ErrorResponse(customMessage)
}

// SuccessResponse is successful operation of any request
func SuccessResponse(message string) gin.H {
	return gin.H{
		"error":   false,
		"message": message,
	}
}

// SuccessResponseWithData will encode data to json or return error
func SuccessResponseWithData(data interface{}) (gin.H, error) {
	return gin.H{
		"error": false,
		"data":  data,
	}, nil
}

// ReportResponse is model response of GetResponse or => /api/report
type ReportResponse struct {
	InstanceID string    `json:"instance_id"`
	Status     int       `json:"status"`
	URL        string    `json:"url"`
	ReportedAt time.Time `json:"reported_at"`
}
