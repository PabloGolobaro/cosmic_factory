package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s service) Get(ctx context.Context, id string) (*model.Order, error) {
	orderUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrInvalidUUID, err)
	}
	order, err := s.Repository.Get(ctx, orderUUID)
	if err != nil {
		return nil, err
	}

	items, err := s.OrderItemRepository.ListByOrder(ctx, orderUUID)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		switch item.PartType {
		case model.PartTypeHull:
			order.HullUUID = item.PartUUID
		case model.PartTypeEngine:
			order.EngineUUID = item.PartUUID
		case model.PartTypeShield:
			order.ShieldUUID = &item.PartUUID
		case model.PartTypeWeapon:
			order.WeaponUUID = &item.PartUUID
		}
	}

	return &order, nil
}
