package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s service) Cancel(ctx context.Context, id string) error {
	orderUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrInvalidUUID, err)
	}

	return s.txManager.Do(ctx, func(txCtx context.Context) error {
		order, err := s.Repository.GetForUpdate(txCtx, orderUUID)
		if err != nil {
			return fmt.Errorf("%w: %w", errs.ErrOrderNotFound, err)
		}

		switch order.Status {
		case model.OrderStatusCancelled:
			return fmt.Errorf("%w: заказ уже отменён", errs.ErrOrderCancelled)
		case model.OrderStatusPaid:
			return fmt.Errorf("%w: заказ уже оплачен", errs.ErrOrderAlreadyPaid)
		}

		uuids := []string{order.HullUUID.String(), order.EngineUUID.String()}
		if sh := uuidPtrToString(order.ShieldUUID); sh != "" {
			uuids = append(uuids, sh)
		}
		if w := uuidPtrToString(order.WeaponUUID); w != "" {
			uuids = append(uuids, w)
		}
		if err = s.InventoryClient.ReleaseParts(txCtx, uuids); err != nil {
			return err
		}

		order.Status = model.OrderStatusCancelled
		return s.Repository.Update(txCtx, order)
	})
}
