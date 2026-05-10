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
	shieldStr := uuidPtrToString(order.ShieldUUID)
	weaponStr := uuidPtrToString(order.WeaponUUID)

	uuids := []string{order.HullUUID.String(), order.EngineUUID.String()}
	if shieldStr != "" {
		uuids = append(uuids, shieldStr)
	}
	if weaponStr != "" {
		uuids = append(uuids, weaponStr)
	}

	parts, err := s.InventoryClient.ListParts(ctx, uuids)
	if err != nil {
		return model.Order{}, err
	}

	partsMap := make(map[string]model.Part, len(parts))
	for _, p := range parts {
		partsMap[p.UUID.String()] = p
	}

	totalPrice, err := calcTotalPrice(partsMap, uuids)
	if err != nil {
		return model.Order{}, err
	}

	if err = s.InventoryClient.ValidateCompatibility(ctx,
		order.HullUUID.String(), order.EngineUUID.String(), shieldStr, weaponStr,
	); err != nil {
		return model.Order{}, err
	}

	if err = s.InventoryClient.ReserveParts(ctx, uuids); err != nil {
		return model.Order{}, err
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
	if shieldStr != "" {
		entries = append(entries, partEntry{*order.ShieldUUID, model.PartTypeShield, partsMap[shieldStr].Price})
	}
	if weaponStr != "" {
		entries = append(entries, partEntry{*order.WeaponUUID, model.PartTypeWeapon, partsMap[weaponStr].Price})
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

func calcTotalPrice(partsMap map[string]model.Part, uuids []string) (int64, error) {
	var total int64
	for _, id := range uuids {
		p, ok := partsMap[id]
		if !ok {
			return 0, fmt.Errorf("%w: %s", errs.ErrPartNotFound, id)
		}
		if p.StockQuantity <= 0 {
			return 0, fmt.Errorf("%w: %s", errs.ErrOutOfStock, p.Name)
		}
		total += p.Price
	}
	return total, nil
}

func uuidPtrToString(u *uuid.UUID) string {
	if u == nil {
		return ""
	}
	return u.String()
}
