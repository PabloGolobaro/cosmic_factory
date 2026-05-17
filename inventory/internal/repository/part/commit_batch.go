package part

import (
	"context"
	"fmt"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
)

func (s *store) CommitBatch(ctx context.Context, parts []entity.Part) error {
	uuids := make([]string, len(parts))
	stocks := make([]int64, len(parts))
	reserved := make([]int, len(parts))
	for i, p := range parts {
		uuids[i] = p.UUID()
		stocks[i] = p.StockQuantity()
		reserved[i] = p.Reserved()
	}

	sql := `UPDATE parts AS p
               SET stock_quantity = batch.stock,
                   reserved       = batch.reserved
              FROM unnest($1::uuid[], $2::bigint[], $3::int[]) AS batch(uuid, stock, reserved)
             WHERE p.uuid = batch.uuid`

	tag, err := s.getter.DefaultTrOrDB(ctx, s.pool).Exec(ctx, sql, uuids, stocks, reserved)
	if err != nil {
		return fmt.Errorf("списать детали: %w", err)
	}
	if tag.RowsAffected() != int64(len(parts)) {
		return fmt.Errorf("%w: обновлено %d из %d", errs.ErrPartNotFound, tag.RowsAffected(), len(parts))
	}
	return nil
}
