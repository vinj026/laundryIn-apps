package utils

import "github.com/gin-gonic/gin"

// SuccessResponse sends a standardized success JSON response.
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status": "success",
		"data":   data,
	})
}

// ErrorResponse sends a standardized error JSON response.
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
