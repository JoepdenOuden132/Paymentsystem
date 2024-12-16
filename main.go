package main

import (
	"main.go/controllers"
	"main.go/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	r := gin.Default()

	r.POST("/payments", controllers.CreatePayment)

	r.Run(":8080")
}
