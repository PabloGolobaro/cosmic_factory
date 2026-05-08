package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
)

type PartType string

const (
	PartTypeUnspecified PartType = "UNSPECIFIED"
	PartTypeHull        PartType = "HULL"
	PartTypeEngine      PartType = "ENGINE"
	PartTypeShield      PartType = "SHIELD"
	PartTypeWeapon      PartType = "WEAPON"
)

// NewPartType создаёт тип детали с валидацией.
func NewPartType(s string) (PartType, error) {
	pt := PartType(s)
	switch pt {
	case PartTypeHull, PartTypeEngine, PartTypeShield, PartTypeWeapon:
		return pt, nil
	default:
		return "", fmt.Errorf("неизвестный тип детали %q: %w", s, errs.ErrInvalidProperties)
	}
}

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64
	PartType      PartType
	StockQuantity int64
	CreatedAt     time.Time
}
