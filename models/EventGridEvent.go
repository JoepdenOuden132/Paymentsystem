package models

import (
	"time"
)

type EventGridEvent struct {
	ID          string      `json:"id"`
	Subject     string      `json:"subject"`
	EventType   string      `json:"eventType"`
	EventTime   time.Time   `json:"eventTime"`
	Data        interface{} `json:"data"`
	DataVersion string      `json:"dataVersion"`
}
