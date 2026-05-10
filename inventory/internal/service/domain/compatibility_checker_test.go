package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
)

func hullPart(strength int) entity.Part {
	props, err := valueobject.NewHullProperties(strength)
	if err != nil {
		panic(err)
	}
	return entity.RestorePart("", "Корпус", "", valueobject.PartTypeHull, 0, 1, 0, props, time.Time{})
}

func enginePart(class valueobject.EngineClass, required int) entity.Part {
	props, err := valueobject.NewEngineProperties(class, required)
	if err != nil {
		panic(err)
	}
	return entity.RestorePart("", "Двигатель", "", valueobject.PartTypeEngine, 0, 1, 0, props, time.Time{})
}

func shieldPart(st valueobject.ShieldType) entity.Part {
	props, err := valueobject.NewShieldProperties(st)
	if err != nil {
		panic(err)
	}
	return entity.RestorePart("", "Щит", "", valueobject.PartTypeShield, 0, 1, 0, props, time.Time{})
}

func weaponPart(wt valueobject.WeaponType) entity.Part {
	props, err := valueobject.NewWeaponProperties(wt)
	if err != nil {
		panic(err)
	}
	return entity.RestorePart("", "Оружие", "", valueobject.PartTypeWeapon, 0, 1, 0, props, time.Time{})
}

func TestCompatibilityChecker_Check(t *testing.T) {
	hull30 := hullPart(30)
	hull50 := hullPart(50)
	engineA40 := enginePart(valueobject.EngineClassA, 40)
	engineB60 := enginePart(valueobject.EngineClassB, 60)
	shieldPlasma := shieldPart(valueobject.ShieldTypePlasma)
	shieldEnergy := shieldPart(valueobject.ShieldTypeEnergy)
	weaponLaser := weaponPart(valueobject.WeaponTypeLaser)
	weaponMissile := weaponPart(valueobject.WeaponTypeMissile)

	checker := NewCompatibilityChecker()

	tests := []struct {
		name    string
		slots   model.ResolvedShipSlots
		wantErr error
	}{
		{
			name:  "корпус и двигатель совместимы, опциональные слоты пусты",
			slots: model.ResolvedShipSlots{Hull: hull50, Engine: engineA40},
		},
		{
			name:    "корпус слишком слабый для двигателя",
			slots:   model.ResolvedShipSlots{Hull: hull30, Engine: engineB60},
			wantErr: errs.ErrIncompatibleParts,
		},
		{
			name:    "плазменный щит несовместим с лазерным оружием",
			slots:   model.ResolvedShipSlots{Hull: hull50, Engine: engineA40, Shield: new(shieldPlasma), Weapon: new(weaponLaser)},
			wantErr: errs.ErrIncompatibleParts,
		},
		{
			name:  "энергетический щит совместим с лазерным оружием",
			slots: model.ResolvedShipSlots{Hull: hull50, Engine: engineA40, Shield: new(shieldEnergy), Weapon: new(weaponLaser)},
		},
		{
			name:  "плазменный щит совместим с ракетным оружием",
			slots: model.ResolvedShipSlots{Hull: hull50, Engine: engineA40, Shield: new(shieldPlasma), Weapon: new(weaponMissile)},
		},
		{
			name:  "все четыре детали полностью совместимы",
			slots: model.ResolvedShipSlots{Hull: hull50, Engine: engineA40, Shield: new(shieldEnergy), Weapon: new(weaponMissile)},
		},
		{
			name:    "все четыре — корпус и двигатель несовместимы",
			slots:   model.ResolvedShipSlots{Hull: hull30, Engine: engineB60, Shield: new(shieldEnergy), Weapon: new(weaponMissile)},
			wantErr: errs.ErrIncompatibleParts,
		},
		{
			name:    "все четыре — щит и оружие несовместимы",
			slots:   model.ResolvedShipSlots{Hull: hull50, Engine: engineA40, Shield: new(shieldPlasma), Weapon: new(weaponLaser)},
			wantErr: errs.ErrIncompatibleParts,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checker.Check(tt.slots)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
