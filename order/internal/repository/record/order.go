package record

import "time"

// OrderRecord модель для БД.
type OrderRecord struct {
	OrderUUID       string
	UserUUID        string
	TotalPrice      int64
	Status          string
	TransactionUUID *string
	PaymentMethod   *string
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}
