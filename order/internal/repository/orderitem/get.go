package orderitem

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/converter"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/record"
)

func (s *repo) Get(ctx context.Context, id uuid.UUID) (model.OrderItem, error) {
	sql := `SELECT * FROM order_items WHERE uuid = $1`

	result := record.OrderItemRecord{}
	err := s.getter.DefaultTrOrDB(ctx, s.pool).QueryRow(ctx, sql, id).Scan(
		&result.UUID, &result.OrderUUID, &result.PartUUID,
		&result.PartType, &result.Price, &result.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.OrderItem{}, fmt.Errorf("%w: %s", errs.ErrOrderItemNotFound, id)
		}
		return model.OrderItem{}, err
	}

	return converter.OrderItemFromRecord(result), nil
}
