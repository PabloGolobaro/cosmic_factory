package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s service) Pay(ctx context.Context, id string, method model.PaymentMethod) (string, error) {
	orderUUID, err := uuid.Parse(id)
	if err != nil {
		return "", fmt.Errorf("%w: %w", errs.ErrInvalidUUID, err)
	}

	var transactionUUID string
	err = s.txManager.Do(ctx, func(txCtx context.Context) error {
		order, err := s.Repository.GetForUpdate(txCtx, orderUUID)
		if err != nil {
			return err
		}
		switch order.Status {
		case model.OrderStatusCancelled:
			return errs.ErrOrderCancelled
		case model.OrderStatusPaid:
			return errs.ErrOrderAlreadyPaid
		}
		transactionUUID, err = s.PaymentClient.PayOrder(txCtx, id, method)
		if err != nil {
			return err
		}
		txUUID, err := uuid.Parse(transactionUUID)
		if err != nil {
			return fmt.Errorf("%w: %w", errs.ErrInvalidUUID, err)
		}
		order.PaymentMethod = method
		order.TransactionUUID = &txUUID
		order.Status = model.OrderStatusPaid
		return s.Repository.Update(txCtx, order)
	})
	return transactionUUID, err
}
