package models

type Payment struct {
	ID            uint    `json:"id"`
	ReservationID uint    `json:"reservation_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	PaymentDate   string  `json:"payment_date"`
}
