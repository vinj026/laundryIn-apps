package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// SuccessResponse sends a standardized success JSON response.
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

// ErrorResponse sends a standardized error JSON response and logs to terminal.
func ErrorResponse(c *gin.Context, statusCode int, message string, errs interface{}) {
	// Log internal server errors to terminal for debugging
	if statusCode >= 500 {
		fmt.Printf("🔴 SERVER ERROR (%d): %s | Details: %+v\n", statusCode, message, errs)
	}

	response := gin.H{
		"status":  "error",
		"message": message,
	}
	if errs != nil {
		// In production, you might want to hide detailed errs from the client,
		// but for now we'll keep them if provided.
		response["errors"] = errs
	}
	c.JSON(statusCode, response)
}
