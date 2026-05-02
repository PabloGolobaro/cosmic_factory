package order

import (
	"context"
	"fmt"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/converter"
)

func (s *repo) Update(ctx context.Context, order model.Order) error {
	rec := converter.OrderToRecord(order)

	sql := `UPDATE orders
	        SET total_price = $2, status = $3, transaction_uuid = $4, payment_method = $5, updated_at = NOW()
	        WHERE uuid = $1`

	cmdTag, err := s.getter.DefaultTrOrDB(ctx, s.pool).Exec(ctx, sql,
		rec.OrderUUID, rec.TotalPrice, rec.Status,
		rec.TransactionUUID, rec.PaymentMethod,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: %s", errs.ErrOrderNotFound, order.OrderUUID)
	}
	return nil
}
