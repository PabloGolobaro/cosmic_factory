package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/converter"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/record"
)

func (s *repo) Create(ctx context.Context, order model.Order) (model.Order, error) {
	order.OrderUUID = uuid.New()
	rec := converter.OrderToRecord(order)

	sql := `INSERT INTO orders (uuid, total_price, status, transaction_uuid, payment_method)
	        VALUES ($1, $2, $3, $4, $5)
	        RETURNING *`

	result := record.OrderRecord{}
	err := s.getter.DefaultTrOrDB(ctx, s.pool).QueryRow(ctx, sql,
		rec.OrderUUID, rec.TotalPrice, rec.Status,
		rec.TransactionUUID, rec.PaymentMethod,
	).Scan(
		&result.OrderUUID, &result.TotalPrice, &result.Status,
		&result.TransactionUUID, &result.PaymentMethod,
		&result.CreatedAt, &result.UpdatedAt,
	)
	if err != nil {
		return model.Order{}, err
	}

	return converter.OrderFromRecord(result), nil
}
