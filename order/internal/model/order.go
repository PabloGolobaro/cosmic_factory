package model

import (
	"time"

	"github.com/google/uuid"
)

// Order представляет заказ на постройку космического корабля.
type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	TotalPrice      int64      // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   PaymentMethod
	Status          string // PENDING_PAYMENT, PAID, CANCELLED
	CreatedAt       time.Time
}
