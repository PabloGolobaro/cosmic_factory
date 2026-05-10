package domain

import (
	"fmt"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
)

type compatibilityChecker struct{}

func NewCompatibilityChecker() *compatibilityChecker {
	return &compatibilityChecker{}
}

func (c *compatibilityChecker) Check(slots entity.ResolvedShipSlots) error {
	hullProps := slots.Hull.Properties()
	engineProps := slots.Engine.Properties()

	hull := hullProps.Hull()
	engine := engineProps.Engine()

	if hull != nil && engine != nil && !hull.CanSupport(engine) {
		return fmt.Errorf(
			"корпус (прочность %d) не выдерживает двигатель класса %s (требует %d): %w",
			hull.Strength(), engine.Class(), engine.RequiredStrength(), errs.ErrIncompatibleParts,
		)
	}

	if slots.Shield != nil && slots.Weapon != nil {
		shieldProps := slots.Shield.Properties()
		weaponProps := slots.Weapon.Properties()

		shield := shieldProps.Shield()
		weapon := weaponProps.Weapon()
		if shield != nil && weapon != nil && shield.ConflictsWith(weapon) {
			return fmt.Errorf("плазменный щит несовместим с лазерным оружием: %w", errs.ErrIncompatibleParts)
		}
	}

	return nil
}
