package orderitem

import (
	"context"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/converter"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/record"
)

func (s *repo) Create(ctx context.Context, item model.OrderItem) (model.OrderItem, error) {
	rec := converter.OrderItemToRecord(item)

	sql := `INSERT INTO order_items (order_uuid, part_uuid, part_type, price)
	        VALUES ($1, $2, $3, $4)
	        RETURNING *`

	result := record.OrderItemRecord{}
	err := s.getter.DefaultTrOrDB(ctx, s.pool).QueryRow(ctx, sql,
		rec.OrderUUID, rec.PartUUID, rec.PartType, rec.Price,
	).Scan(
		&result.UUID, &result.OrderUUID, &result.PartUUID,
		&result.PartType, &result.Price, &result.CreatedAt,
	)
	if err != nil {
		return model.OrderItem{}, err
	}

	return converter.OrderItemFromRecord(result), nil
}
