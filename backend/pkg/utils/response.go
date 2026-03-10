package utils

import "github.com/gin-gonic/gin"

// SuccessResponse sends a standardized success JSON response.
// Format: {"status": "success", "message": "...", "data": ...}
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

// ErrorResponse sends a standardized error JSON response.
// Format: {"status": "error", "message": "...", "errors": ...}
func ErrorResponse(c *gin.Context, statusCode int, message string, errs interface{}) {
	response := gin.H{
		"status":  "error",
		"message": message,
	}
	if errs != nil {
		response["errors"] = errs
	}
	c.JSON(statusCode, response)
}
