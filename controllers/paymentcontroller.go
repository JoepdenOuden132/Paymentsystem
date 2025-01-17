package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"main.go/models"
)

type PaymentManager struct {
	payments []models.Payment
	mu       sync.Mutex
	nextID   uint
}

func NewPaymentManager() *PaymentManager {
	return &PaymentManager{
		payments: []models.Payment{},
		nextID:   1,
	}
}

func (pm *PaymentManager) CreatePayment(newPayment models.Payment) uint {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	newPayment.ID = pm.nextID
	pm.nextID++
	pm.payments = append(pm.payments, newPayment)
	log.Printf("Payment created with ID: %d", newPayment.ID)
	return newPayment.ID
}

func (pm *PaymentManager) UpdatePaymentStatus(paymentID uint, status string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for i, payment := range pm.payments {
		if payment.ID == paymentID {
			pm.payments[i].Status = status
			log.Printf("Updated payment ID %d with status: %s", paymentID, status)
			break
		}
	}
}

var paymentManager = NewPaymentManager()

func CreatePayment(c *gin.Context) {
	var newPayment models.Payment

	// Keurt de binnenkomende JSON goed
	if err := c.ShouldBindJSON(&newPayment); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Stelt de initiële waarden in
	newPayment.Status = "pending"
	newPayment.PaymentDate = time.Now().Format("2006-01-02 15:04:05")

	paymentID := paymentManager.CreatePayment(newPayment)

	// Begint versturen naar Event Grid
	go func() {
		if err := sendEventToEventGrid(newPayment); err != nil {
			log.Printf("Error sending event to Event Grid: %v", err)
		}
	}()

	go simulatePayment(paymentID)

	c.JSON(http.StatusCreated, gin.H{"message": "Payment is being processed", "payment_id": paymentID})
}

func simulatePayment(paymentID uint) {
	time.Sleep(5 * time.Second)
	statuses := []string{"success", "failed"}
	rand.Seed(time.Now().UnixNano())
	randomStatus := statuses[rand.Intn(len(statuses))]
	log.Printf("Simulating payment ID %d with status: %s", paymentID, randomStatus)
	paymentManager.UpdatePaymentStatus(paymentID, randomStatus)
}

func sendEventToEventGrid(payment models.Payment) error {
	event := []map[string]interface{}{
		{
			"id":        "20",
			"eventType": "Payment.Created",
			"subject":   "new/payment",
			"eventTime": time.Now().Format(time.RFC3339),
			"data": map[string]interface{}{
				"paymentId":   payment.ID,
				"amount":      payment.Amount,
				"status":      payment.Status,
				"paymentDate": payment.PaymentDate,
			},
			"dataVersion": "1.0",
		},
	}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	eventGridEndpoint := os.Getenv("EVENT_GRID_ENDPOINT")
	eventGridKey := os.Getenv("EVENT_GRID_KEY")

	log.Println("Event Grid Endpoint:", eventGridEndpoint)
	log.Println("Event Grid Key:", eventGridKey)

	retryCount := 1
	for i := 0; i < retryCount; i++ {
		req, err := http.NewRequest("POST", eventGridEndpoint, bytes.NewBuffer(eventJSON))
		if err != nil {
			log.Printf("Attempt %d: failed to create request: %v", i+1, err)
			return fmt.Errorf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("aeg-sas-key", eventGridKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			//	log.Printf("Attempt %d: failed to send request: %v", i+1, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
			log.Printf("Successfully sent event to Event Grid on attempt %d", i+1)
			return nil
		}
		log.Printf("Successfully sent event to Event Grid on attempt %d", i+1)
	}

	// log.Printf("Failed to send event after %d attempts", retryCount)
	return nil
}
