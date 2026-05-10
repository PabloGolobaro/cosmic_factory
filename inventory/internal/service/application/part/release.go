package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
)

func (s *service) ReleaseParts(ctx context.Context, uuids []string) error {
	for _, id := range uuids {
		if _, err := uuid.Parse(id); err != nil {
			return fmt.Errorf("%w: %w", errs.ErrInvalidUUID, err)
		}
	}

	return s.txManager.Do(ctx, func(ctx context.Context) error {
		parts, err := s.PartRepository.GetBatch(ctx, valueobject.PartFilter{UUIDs: uuids})
		if err != nil {
			return fmt.Errorf("получить детали: %w", err)
		}

		for i := range parts {
			if err = parts[i].Release(1); err != nil {
				return err
			}
		}

		return s.PartRepository.UpdateReservedBatch(ctx, parts)
	})
}
