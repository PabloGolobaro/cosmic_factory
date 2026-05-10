package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
)

func (s *service) ValidateCompatibility(ctx context.Context, slots model.ShipSlots) error {
	resolved, err := s.resolveShipSlots(ctx, slots)
	if err != nil {
		return err
	}

	return s.compatibilityChecker.Check(resolved)
}

func (s *service) resolveShipSlots(ctx context.Context, slots model.ShipSlots) (model.ResolvedShipSlots, error) {
	type slotSpec struct {
		uuid     string
		expected valueobject.PartType
	}

	specs := []slotSpec{
		{slots.HullUUID, valueobject.PartTypeHull},
		{slots.EngineUUID, valueobject.PartTypeEngine},
	}
	if slots.ShieldUUID != "" {
		specs = append(specs, slotSpec{slots.ShieldUUID, valueobject.PartTypeShield})
	}
	if slots.WeaponUUID != "" {
		specs = append(specs, slotSpec{slots.WeaponUUID, valueobject.PartTypeWeapon})
	}

	uuids := make([]string, 0, len(specs))
	for _, spec := range specs {
		if _, err := uuid.Parse(spec.uuid); err != nil {
			return model.ResolvedShipSlots{}, fmt.Errorf("%w: %s", errs.ErrInvalidUUID, spec.uuid)
		}
		uuids = append(uuids, spec.uuid)
	}

	fetched, err := s.PartRepository.GetBatch(ctx, model.PartFilter{UUIDs: uuids})
	if err != nil {
		return model.ResolvedShipSlots{}, fmt.Errorf("получить детали: %w", err)
	}

	byUUID := make(map[string]entity.Part, len(fetched))
	for _, p := range fetched {
		byUUID[p.UUID()] = p
	}

	get := func(id string, expected valueobject.PartType) (entity.Part, error) {
		p, ok := byUUID[id]
		if !ok {
			return entity.Part{}, fmt.Errorf("деталь %s: %w", id, errs.ErrPartNotFound)
		}
		if p.PartType() != expected {
			return entity.Part{}, fmt.Errorf(
				"слот %s: ожидается %s, получен %s: %w",
				id, expected, p.PartType(), errs.ErrPartTypeMismatch,
			)
		}
		return p, nil
	}

	hull, err := get(slots.HullUUID, valueobject.PartTypeHull)
	if err != nil {
		return model.ResolvedShipSlots{}, err
	}
	engine, err := get(slots.EngineUUID, valueobject.PartTypeEngine)
	if err != nil {
		return model.ResolvedShipSlots{}, err
	}

	resolved := model.ResolvedShipSlots{Hull: hull, Engine: engine}

	if slots.ShieldUUID != "" {
		shield, err := get(slots.ShieldUUID, valueobject.PartTypeShield)
		if err != nil {
			return model.ResolvedShipSlots{}, err
		}
		resolved.Shield = &shield
	}
	if slots.WeaponUUID != "" {
		weapon, err := get(slots.WeaponUUID, valueobject.PartTypeWeapon)
		if err != nil {
			return model.ResolvedShipSlots{}, err
		}
		resolved.Weapon = &weapon
	}

	return resolved, nil
}
