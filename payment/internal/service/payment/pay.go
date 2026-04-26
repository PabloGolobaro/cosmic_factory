package payment

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/payment/internal/errors"
)

func (p paymentService) Pay(ctx context.Context, id, paymentMethod string) (string, error) {
	if paymentMethod == "" {
		return "", errs.ErrInvalidPaymentMethod
	}

	if _, err := uuid.Parse(id); err != nil {
		return "", fmt.Errorf("%w: %w", errs.ErrInvalidUUID, err)
	}

	txUUID := uuid.NewString()

	slog.Info("оплата прошла успешно", "order_uuid", id, "transaction_uuid", txUUID)

	return txUUID, nil
}
