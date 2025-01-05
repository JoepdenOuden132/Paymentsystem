package main

import (
	"main.go/controllers"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.POST("/payments", controllers.CreatePayment)

	r.Run(":8080")
}
