package domain

import (
	"fmt"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
)

type compatibilityChecker struct{}

func NewCompatibilityChecker() *compatibilityChecker {
	return &compatibilityChecker{}
}

func (c *compatibilityChecker) Check(parts []entity.Part) error {
	var hull *valueobject.HullProperties
	var engine *valueobject.EngineProperties
	var shield *valueobject.ShieldProperties
	var weapon *valueobject.WeaponProperties

	for _, p := range parts {
		props := p.Properties()
		if h := props.Hull(); h != nil {
			hull = h
		}
		if e := props.Engine(); e != nil {
			engine = e
		}
		if s := props.Shield(); s != nil {
			shield = s
		}
		if w := props.Weapon(); w != nil {
			weapon = w
		}
	}

	if hull != nil && engine != nil && !hull.CanSupport(engine) {
		return fmt.Errorf(
			"корпус (прочность %d) не выдерживает двигатель класса %s (требует %d): %w",
			hull.Strength(), engine.Class(), engine.RequiredStrength(), errs.ErrIncompatibleParts,
		)
	}

	if shield != nil && weapon != nil && shield.ConflictsWith(weapon) {
		return fmt.Errorf("плазменный щит несовместим с лазерным оружием: %w", errs.ErrIncompatibleParts)
	}

	return nil
}
