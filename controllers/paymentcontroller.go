package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"main.go/models"
)

// Simpele slice om betalingen in op te slaan (in-memory opslag)
var payments []models.Payment
var mu sync.Mutex // Mutex om concurrente toegang te beheren

var nextID uint = 1 // Simuleer auto-increment ID

func CreatePayment(c *gin.Context) {
	var newPayment models.Payment

	// Valideer de binnenkomende JSON
	if err := c.ShouldBindJSON(&newPayment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Stel de initiële waarden in
	newPayment.ID = getNextID()
	newPayment.Status = "pending"
	newPayment.PaymentDate = time.Now().Format("2006-01-02 15:04:05")

	// Voeg de nieuwe betaling toe aan de slice
	mu.Lock()
	payments = append(payments, newPayment)
	mu.Unlock()

	go func() {
		if err := sendEventToEventGrid(newPayment); err != nil {
			fmt.Printf("Error sending event to Event Grid: %v\n", err)
		}
	}()

	go simulatePayment(newPayment.ID)

	c.JSON(http.StatusCreated, gin.H{"message": "Payment is being processed", "payment_id": newPayment.ID})
}

func simulatePayment(paymentID uint) {
	time.Sleep(5 * time.Second)
	statuses := []string{"success", "failed"}
	rand.Seed(time.Now().UnixNano())
	randomStatus := statuses[rand.Intn(len(statuses))]

	// Update de status van de betaling
	mu.Lock()
	for i, payment := range payments {
		if payment.ID == paymentID {
			payments[i].Status = randomStatus
			break
		}
	}
	mu.Unlock()
}

// Hulpfunctie om het volgende ID te genereren
func getNextID() uint {
	mu.Lock()
	defer mu.Unlock()
	nextID++
	return nextID - 1
}

func sendEventToEventGrid(payment models.Payment) error {
	event := []map[string]interface{}{
		{
			"id":          fmt.Sprintf("%d", payment.ID),
			"eventType":   "Payment.Created",
			"subject":     "new/payment",
			"eventTime":   time.Now().Format(time.RFC3339),
			"data":        payment,
			"dataVersion": "1.0",
		},
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	eventGridEndpoint := os.Getenv("EVENT_GRID_ENDPOINT") // Zet dit in je container als environment variable
	eventGridKey := os.Getenv("EVENT_GRID_KEY")           // Zet dit in je container als environment variable

	req, err := http.NewRequest("POST", eventGridEndpoint, bytes.NewBuffer(eventJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("aeg-sas-key", eventGridKey) // Event Grid authenticatie

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("received non-OK response: %v", resp.Status)
	}

	return nil
}
