package part

import (
	"context"
	"fmt"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
)

func (s *store) UpdateReservedBatch(ctx context.Context, parts []entity.Part) error {
	uuids := make([]string, len(parts))
	reserved := make([]int, len(parts))
	for i, p := range parts {
		uuids[i] = p.UUID()
		reserved[i] = p.Reserved()
	}

	sql := `UPDATE parts AS p
	   SET reserved = batch.reserved
	  FROM unnest($1::uuid[], $2::int[]) AS batch(uuid, reserved)
	 WHERE p.uuid = batch.uuid`

	_, err := s.getter.DefaultTrOrDB(ctx, s.pool).Exec(ctx, sql, uuids, reserved)
	if err != nil {
		return fmt.Errorf("обновить резерв: %w", err)
	}

	return nil
}
