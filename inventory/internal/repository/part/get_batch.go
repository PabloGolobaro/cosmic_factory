package part

import (
	"context"
	"fmt"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/converter"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/record"
)

const selectPartCols = `SELECT uuid, name, description, part_type, price, stock_quantity, reserved, properties, created_at FROM parts`

func (s *store) GetBatch(ctx context.Context, filter valueobject.PartFilter) ([]entity.Part, error) {
	var (
		sqlStr string
		args   []any
	)

	switch {
	case len(filter.UUIDs) > 0:
		sqlStr = selectPartCols + ` WHERE uuid = ANY($1)`
		args = []any{filter.UUIDs}
	case filter.PartType != valueobject.PartTypeUnspecified:
		sqlStr = selectPartCols + ` WHERE part_type = $1`
		args = []any{string(filter.PartType)}
	default:
		sqlStr = selectPartCols
	}

	rows, err := s.getter.DefaultTrOrDB(ctx, s.pool).Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scanRow := func() (entity.Part, error) {
		r := record.PartRecord{}
		if err = rows.Scan(&r.UUID, &r.Name, &r.Description,
			&r.PartType, &r.Price, &r.StockQuantity, &r.Reserved, &r.Properties, &r.CreatedAt); err != nil {
			return entity.Part{}, err
		}
		return converter.PartFromRecord(r)
	}

	if len(filter.UUIDs) > 0 {
		found := make(map[string]entity.Part, len(filter.UUIDs))
		for rows.Next() {
			p, err := scanRow()
			if err != nil {
				return nil, fmt.Errorf("ошибка конвертации: %w", err)
			}
			found[p.UUID()] = p
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}

		parts := make([]entity.Part, 0, len(filter.UUIDs))
		for _, id := range filter.UUIDs {
			p, ok := found[id]
			if !ok {
				return nil, fmt.Errorf("%w: %s", errs.ErrPartNotFound, id)
			}
			parts = append(parts, p)
		}
		return parts, nil
	}

	var parts []entity.Part
	for rows.Next() {
		p, err := scanRow()
		if err != nil {
			return nil, fmt.Errorf("ошибка конвертации: %w", err)
		}
		parts = append(parts, p)
	}
	return parts, rows.Err()
}
