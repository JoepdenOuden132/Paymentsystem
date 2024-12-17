package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"main.go/models" // pas dit aan naar jouw package

	"github.com/google/uuid" // Zorg dat je deze library installeert
)

// SendEventToEventGrid stuurt een event naar Azure Event Grid
func SendEventToEventGrid(eventType string, payment models.Payment) error {
	event := models.EventGridEvent{
		ID:          uuid.New().String(),
		Subject:     fmt.Sprintf("payment/%d", payment.ID),
		EventType:   eventType,
		EventTime:   time.Now(),
		Data:        payment,
		DataVersion: "1.0",
	}

	// Zet de EventGridEvent om naar JSON
	eventData, err := json.Marshal([]models.EventGridEvent{event}) // array van events
	if err != nil {
		log.Printf("Fout bij het omzetten van event naar JSON: %v", err)
		return err
	}

	// Haal de Event Grid-configuratie op uit de omgeving
	eventGridURL := os.Getenv("EVENT_GRID_TOPIC_ENDPOINT")
	eventGridKey := os.Getenv("EVENT_GRID_ACCESS_KEY")

	// Maak een POST-aanvraag naar Event Grid
	req, err := http.NewRequest("POST", eventGridURL, bytes.NewBuffer(eventData))
	if err != nil {
		log.Printf("Fout bij het maken van de HTTP-aanvraag: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("aeg-sas-key", eventGridKey) // Event Grid authenticatie

	// Verstuur het verzoek
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Fout bij het versturen van de HTTP-aanvraag: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		log.Printf("Event Grid stuurde een onverwachte statuscode terug: %d", resp.StatusCode)
		return fmt.Errorf("Event Grid stuurde een onverwachte statuscode terug: %d", resp.StatusCode)
	}

	log.Printf("Event succesvol verstuurd naar Event Grid! Statuscode: %d", resp.StatusCode)
	return nil
}
