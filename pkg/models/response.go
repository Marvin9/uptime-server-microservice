package models

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse is response of any request, if any error
func ErrorResponse(err string) gin.H {
	return gin.H{
		"error":   true,
		"message": err,
	}
}

// SuccessResponse is successful operation of any request
func SuccessResponse(message string) gin.H {
	return gin.H{
		"error":   false,
		"message": message,
	}
}
