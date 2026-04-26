package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s service) Pay(ctx context.Context, id string) (uuid.UUID, error) {
	orderUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %w", errs.ErrInvalidUUID, err)
	}

	order, err := s.Repository.Get(ctx, orderUUID)
	if err != nil {
		return uuid.Nil, err
	}

	switch order.Status {
	case model.OrderStatusCancelled:
		return uuid.Nil, errs.ErrOrderCancelled
	case model.OrderStatusPaid:
		return uuid.Nil, errs.ErrOrderAlreadyPaid
	}

	if err = s.PaymentClient.PayOrder(ctx, id, order.PaymentMethod); err != nil {
		return uuid.Nil, err
	}

	txUUID := uuid.New()
	order.TransactionUUID = &txUUID
	order.Status = model.OrderStatusPaid
	return txUUID, s.Repository.Update(ctx, order)
}
