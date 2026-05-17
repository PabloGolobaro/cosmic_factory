package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s service) CommitShipParts(ctx context.Context, orderUUID string) error {
	id, err := uuid.Parse(orderUUID)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrInvalidUUID, err)
	}

	return s.txManager.Do(ctx, func(txCtx context.Context) error {
		order, err := s.Repository.GetForUpdate(txCtx, id)
		if err != nil {
			return fmt.Errorf("%w: %w", errs.ErrOrderNotFound, err)
		}

		if order.Status == model.OrderStatusAssembled {
			return nil
		}

		uuids := []string{order.HullUUID.String(), order.EngineUUID.String()}
		if sh := uuidPtrToString(order.ShieldUUID); sh != "" {
			uuids = append(uuids, sh)
		}
		if w := uuidPtrToString(order.WeaponUUID); w != "" {
			uuids = append(uuids, w)
		}

		if err = s.InventoryClient.CommitParts(txCtx, uuids); err != nil {
			return err
		}

		order.Status = model.OrderStatusAssembled
		return s.Repository.Update(txCtx, order)
	})
}
