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

	order, err := s.Repository.Get(ctx, orderUUID)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrOrderNotFound, err)
	}

	switch order.Status {
	case model.OrderStatusCancelled:
		return fmt.Errorf("%w: заказ уже отменён", errs.ErrOrderCancelled)
	case model.OrderStatusPaid:
		return fmt.Errorf("%w: заказ уже оплачен", errs.ErrOrderAlreadyPaid)
	}

	order.Status = model.OrderStatusCancelled
	return s.Repository.Update(ctx, order)
}
