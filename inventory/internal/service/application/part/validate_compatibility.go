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

func collectUUIDs(ids ...string) ([]string, error) {
	var result []string
	for _, id := range ids {
		if id == "" {
			continue
		}
		if _, err := uuid.Parse(id); err != nil {
			return nil, fmt.Errorf("%w: %s", errs.ErrInvalidUUID, id)
		}
		result = append(result, id)
	}
	return result, nil
}

func (s *service) resolveShipSlots(ctx context.Context, slots model.ShipSlots) (model.ResolvedShipSlots, error) {
	if slots.HullUUID == "" {
		return model.ResolvedShipSlots{}, fmt.Errorf("%w: hull_uuid обязателен", errs.ErrInvalidUUID)
	}
	if slots.EngineUUID == "" {
		return model.ResolvedShipSlots{}, fmt.Errorf("%w: engine_uuid обязателен", errs.ErrInvalidUUID)
	}

	uuids, err := collectUUIDs(slots.HullUUID, slots.EngineUUID, slots.ShieldUUID, slots.WeaponUUID)
	if err != nil {
		return model.ResolvedShipSlots{}, err
	}

	fetched, err := s.PartRepository.GetBatch(ctx, model.PartFilter{UUIDs: uuids})
	if err != nil {
		return model.ResolvedShipSlots{}, fmt.Errorf("получить детали: %w", err)
	}

	byUUID := make(map[string]entity.Part, len(fetched))
	for _, p := range fetched {
		byUUID[p.UUID()] = p
	}

	get := func(id string) (entity.Part, error) {
		p, ok := byUUID[id]
		if !ok {
			return entity.Part{}, fmt.Errorf("деталь %s: %w", id, errs.ErrPartNotFound)
		}
		return p, nil
	}

	var resolved model.ResolvedShipSlots

	hull, err := get(slots.HullUUID)
	if err != nil {
		return model.ResolvedShipSlots{}, err
	}
	if hull.PartType() != valueobject.PartTypeHull {
		return model.ResolvedShipSlots{}, fmt.Errorf("%w: слот hull ожидает HULL, получен %s", errs.ErrPartTypeMismatch, hull.PartType())
	}
	resolved.Hull = hull

	engine, err := get(slots.EngineUUID)
	if err != nil {
		return model.ResolvedShipSlots{}, err
	}
	if engine.PartType() != valueobject.PartTypeEngine {
		return model.ResolvedShipSlots{}, fmt.Errorf("%w: слот engine ожидает ENGINE, получен %s", errs.ErrPartTypeMismatch, engine.PartType())
	}
	resolved.Engine = engine

	if slots.ShieldUUID != "" {
		shield, err := get(slots.ShieldUUID)
		if err != nil {
			return model.ResolvedShipSlots{}, err
		}
		if shield.PartType() != valueobject.PartTypeShield {
			return model.ResolvedShipSlots{}, fmt.Errorf("%w: слот shield ожидает SHIELD, получен %s", errs.ErrPartTypeMismatch, shield.PartType())
		}
		resolved.Shield = &shield
	}
	if slots.WeaponUUID != "" {
		weapon, err := get(slots.WeaponUUID)
		if err != nil {
			return model.ResolvedShipSlots{}, err
		}
		if weapon.PartType() != valueobject.PartTypeWeapon {
			return model.ResolvedShipSlots{}, fmt.Errorf("%w: слот weapon ожидает WEAPON, получен %s", errs.ErrPartTypeMismatch, weapon.PartType())
		}
		resolved.Weapon = &weapon
	}

	return resolved, nil
}
