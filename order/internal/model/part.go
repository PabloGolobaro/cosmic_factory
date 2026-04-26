package model

import (
	"time"

	"github.com/google/uuid"
)

type PartType int

const (
	PartTypeUnspecified PartType = iota
	PartTypeHull
	PartTypeEngine
	PartTypeShield
	PartTypeWeapon
)

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64
	PartType      PartType
	StockQuantity int64
	CreatedAt     time.Time
}
