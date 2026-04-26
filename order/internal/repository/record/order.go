package record

import "time"

// OrderRecord модель для БД.
type OrderRecord struct {
	OrderUUID       string
	HullUUID        string
	EngineUUID      string
	ShieldUUID      *string // опциональный
	WeaponUUID      *string // опциональный
	TotalPrice      int64   // в копейках
	TransactionUUID *string
	PaymentMethod   *string
	Status          string // PENDING_PAYMENT, PAID, CANCELLED
	CreatedAt       time.Time
}
