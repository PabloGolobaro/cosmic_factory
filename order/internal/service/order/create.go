package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s service) Create(ctx context.Context, order model.Order) (model.Order, error) {
	uuids := []string{order.HullUUID.String(), order.EngineUUID.String()}
	if order.ShieldUUID != nil {
		uuids = append(uuids, order.ShieldUUID.String())
	}
	if order.WeaponUUID != nil {
		uuids = append(uuids, order.WeaponUUID.String())
	}

	parts, err := s.InventoryClient.ListParts(ctx, uuids)
	if err != nil {
		return model.Order{}, err
	}

	partsMap := make(map[string]model.Part, len(parts))
	for _, p := range parts {
		partsMap[p.UUID.String()] = p
	}

	var totalPrice int64
	for _, id := range uuids {
		p, ok := partsMap[id]
		if !ok {
			return model.Order{}, fmt.Errorf("%w: %s", errs.ErrPartNotFound, id)
		}
		if p.StockQuantity <= 0 {
			return model.Order{}, fmt.Errorf("%w: %s", errs.ErrOutOfStock, p.Name)
		}
		totalPrice += p.Price
	}

	order.TotalPrice = totalPrice
	order.Status = model.OrderStatusPendingPayment
	order.CreatedAt = time.Now()

	type partEntry struct {
		uuid     uuid.UUID
		partType model.PartType
		price    int64
	}
	entries := []partEntry{
		{order.HullUUID, model.PartTypeHull, partsMap[order.HullUUID.String()].Price},
		{order.EngineUUID, model.PartTypeEngine, partsMap[order.EngineUUID.String()].Price},
	}
	if order.ShieldUUID != nil {
		entries = append(entries, partEntry{*order.ShieldUUID, model.PartTypeShield, partsMap[order.ShieldUUID.String()].Price})
	}
	if order.WeaponUUID != nil {
		entries = append(entries, partEntry{*order.WeaponUUID, model.PartTypeWeapon, partsMap[order.WeaponUUID.String()].Price})
	}

	var createdOrder model.Order
	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		var txErr error
		createdOrder, txErr = s.Repository.Create(ctx, order)
		if txErr != nil {
			return txErr
		}
		createdOrder.HullUUID = order.HullUUID
		createdOrder.EngineUUID = order.EngineUUID
		createdOrder.ShieldUUID = order.ShieldUUID
		createdOrder.WeaponUUID = order.WeaponUUID
		for _, e := range entries {
			item := model.OrderItem{
				OrderUUID: createdOrder.OrderUUID,
				PartUUID:  e.uuid,
				PartType:  e.partType,
				Price:     e.price,
			}
			if _, txErr = s.OrderItemRepository.Create(ctx, item); txErr != nil {
				return txErr
			}
		}
		return nil
	})
	if err != nil {
		return model.Order{}, err
	}
	return createdOrder, nil
}
