package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode("debug")
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Admin API server",
		})
	})

	addr := fmt.Sprintf(":%d", 8081)
	log.Printf("Admin server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start admin server: %v", err)
	}
}
