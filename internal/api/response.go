package api

import "github.com/gin-gonic/gin"

func respondError(c *gin.Context, status int, code string, message string) {
	if message == "" {
		message = code
	}
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
