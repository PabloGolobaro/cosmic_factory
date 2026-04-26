package order

import (
	"context"
	"fmt"
	"time"

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

	return s.Repository.Create(ctx, order)
}
