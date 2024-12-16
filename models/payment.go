package models

type Payment struct {
	ID            uint    `json:"id" gorm:"primaryKey"`
	ReservationID uint    `json:"reservation_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Currency      string  `json:"currency" binding:"required"`
	Status        string  `json:"status"`
	PaymentDate   string  `json:"payment_date"`
}
