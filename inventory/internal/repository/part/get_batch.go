package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/converter"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/record"
)

func (s *store) GetBatch(ctx context.Context, ids []uuid.UUID) ([]model.Part, error) {
	sql := `SELECT * FROM parts WHERE uuid = ANY($1)`

	rows, err := s.getter.DefaultTrOrDB(ctx, s.pool).Query(ctx, sql, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	found := make(map[uuid.UUID]model.Part, len(ids))
	for rows.Next() {
		r := record.PartRecord{}
		if err = rows.Scan(&r.UUID, &r.Name, &r.Description,
			&r.PartType, &r.Price, &r.StockQuantity, &r.CreatedAt); err != nil {
			return nil, err
		}
		p := converter.PartFromRecord(r)
		found[p.UUID] = p
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	parts := make([]model.Part, 0, len(ids))
	for _, id := range ids {
		p, ok := found[id]
		if !ok {
			return nil, fmt.Errorf("%w: %s", errs.ErrPartNotFound, id)
		}
		parts = append(parts, p)
	}

	return parts, nil
}
