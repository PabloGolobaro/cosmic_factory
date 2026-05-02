package record

import (
	"time"
)

type PartRecord struct {
	UUID          string
	Name          string
	Description   string
	PartType      string
	Price         int64
	StockQuantity int64
	CreatedAt     time.Time
}
