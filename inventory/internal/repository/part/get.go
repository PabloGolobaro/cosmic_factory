package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/converter"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/record"
)

func (s *store) Get(ctx context.Context, id uuid.UUID) (entity.Part, error) {
	sql := selectPartCols + ` WHERE uuid = $1`

	result := record.PartRecord{}
	err := s.getter.DefaultTrOrDB(ctx, s.pool).QueryRow(ctx, sql, id).Scan(
		&result.UUID, &result.Name, &result.Description,
		&result.PartType, &result.Price, &result.StockQuantity, &result.Reserved, &result.Properties, &result.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Part{}, fmt.Errorf("%w: %s", errs.ErrPartNotFound, id)
		}
		return entity.Part{}, err
	}

	p, err := converter.PartFromRecord(result)
	if err != nil {
		return entity.Part{}, fmt.Errorf("ошибка конвертации записи: %w", err)
	}

	return p, nil
}
