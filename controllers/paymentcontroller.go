package controllers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"main.go/database"
	"main.go/models"
)

func CreatePayment(c *gin.Context) {
	var newPayment models.Payment

	if err := c.ShouldBindJSON(&newPayment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPayment.Status = "pending"
	newPayment.PaymentDate = time.Now().Format("2000-01-01 15:00:00")

	if err := database.DB.Create(&newPayment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	go simulatePayment(newPayment.ID)

	c.JSON(http.StatusCreated, gin.H{"message": "Payment is being processed", "payment_id": newPayment.ID})
}

func simulatePayment(paymentID uint) {
	time.Sleep(5 * time.Second)
	statuses := []string{"success", "failed"}
	rand.Seed(time.Now().UnixNano())
	randomStatus := statuses[rand.Intn(len(statuses))]

	database.DB.Model(&models.Payment{}).Where("id = ?", paymentID).Update("status", randomStatus)
}
