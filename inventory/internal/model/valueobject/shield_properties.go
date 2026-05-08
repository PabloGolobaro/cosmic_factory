package valueobject

import (
	"fmt"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
)

// ShieldType — тип щита.
type ShieldType string

const (
	// ShieldTypeEnergy — энергетический щит.
	ShieldTypeEnergy ShieldType = "energy"
	// ShieldTypePlasma — плазменный щит.
	ShieldTypePlasma ShieldType = "plasma"
)

// NewShieldType создаёт ShieldType с валидацией допустимых значений.
func NewShieldType(s string) (ShieldType, error) {
	st := ShieldType(s)

	switch st {
	case ShieldTypeEnergy, ShieldTypePlasma:
		return st, nil
	default:
		return "", fmt.Errorf("недопустимый тип щита %q: %w", s, errs.ErrInvalidProperties)
	}
}

// ShieldProperties — свойства щита (Value Object).
type ShieldProperties struct {
	shieldType ShieldType
}

func (s *ShieldProperties) ShieldType() ShieldType { return s.shieldType }

func (s *ShieldProperties) ConflictsWith(w *WeaponProperties) bool {
	return s.shieldType == ShieldTypePlasma && w.WeaponType() == WeaponTypeLaser
}

// NewShieldProperties создаёт свойства щита.
func NewShieldProperties(shieldType ShieldType) (PartProperties, error) {
	switch shieldType {
	case ShieldTypeEnergy, ShieldTypePlasma:
	default:
		return PartProperties{}, fmt.Errorf("недопустимый тип щита %q: %w", shieldType, errs.ErrInvalidProperties)
	}

	return PartProperties{
		shield: &ShieldProperties{shieldType: shieldType},
	}, nil
}
