package main

import (
	"main.go/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Routes
	r.POST("/payments", controllers.CreatePayment)

	// Start server
	r.Run(":8080")
}
