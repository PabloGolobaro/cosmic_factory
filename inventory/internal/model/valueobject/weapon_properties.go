package valueobject

import (
	"fmt"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
)

// WeaponType — тип оружия.
type WeaponType string

const (
	// WeaponTypeLaser — лазерное оружие.
	WeaponTypeLaser WeaponType = "laser"
	// WeaponTypeMissile — ракетное оружие.
	WeaponTypeMissile WeaponType = "missile"
)

// NewWeaponType создаёт WeaponType с валидацией допустимых значений.
func NewWeaponType(s string) (WeaponType, error) {
	wt := WeaponType(s)

	switch wt {
	case WeaponTypeLaser, WeaponTypeMissile:
		return wt, nil
	default:
		return "", fmt.Errorf("недопустимый тип оружия %q: %w", s, errs.ErrInvalidProperties)
	}
}

// WeaponProperties — свойства оружия (Value Object).
type WeaponProperties struct {
	weaponType WeaponType
}

func (w *WeaponProperties) WeaponType() WeaponType { return w.weaponType }

// NewWeaponProperties создаёт свойства оружия.
func NewWeaponProperties(weaponType WeaponType) (PartProperties, error) {
	switch weaponType {
	case WeaponTypeLaser, WeaponTypeMissile:
	default:
		return PartProperties{}, fmt.Errorf("недопустимый тип оружия %q: %w", weaponType, errs.ErrInvalidProperties)
	}

	return PartProperties{
		weapon: &WeaponProperties{weaponType: weaponType},
	}, nil
}
