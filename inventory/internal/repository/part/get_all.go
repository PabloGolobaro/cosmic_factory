package part

import (
	"context"
	"fmt"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/converter"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/record"
)

func (s *store) GetAll(ctx context.Context) ([]model.Part, error) {
	sql := `SELECT * FROM parts`

	rows, err := s.getter.DefaultTrOrDB(ctx, s.pool).Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []model.Part
	for rows.Next() {
		r := record.PartRecord{}
		if err = rows.Scan(&r.UUID, &r.Name, &r.Description,
			&r.PartType, &r.Price, &r.StockQuantity, &r.CreatedAt); err != nil {
			return nil, err
		}
		p, err := converter.PartFromRecord(r)
		if err != nil {
			return nil, fmt.Errorf("ошибка конвертации записи: %w", err)
		}
		parts = append(parts, p)
	}

	return parts, rows.Err()
}
