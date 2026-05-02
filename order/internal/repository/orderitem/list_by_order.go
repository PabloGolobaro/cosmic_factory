package orderitem

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/converter"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/record"
)

func (s *repo) ListByOrder(ctx context.Context, orderUUID uuid.UUID) ([]model.OrderItem, error) {
	sql := `SELECT * FROM order_items WHERE order_uuid = $1`

	rows, err := s.getter.DefaultTrOrDB(ctx, s.pool).Query(ctx, sql, orderUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.OrderItem
	for rows.Next() {
		r := record.OrderItemRecord{}
		if err = rows.Scan(&r.UUID, &r.OrderUUID, &r.PartUUID,
			&r.PartType, &r.Price, &r.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, converter.OrderItemFromRecord(r))
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
