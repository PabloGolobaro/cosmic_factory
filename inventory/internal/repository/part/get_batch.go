package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/converter"
)

func (s *store) GetBatch(_ context.Context, ids []uuid.UUID) ([]model.Part, error) {
	parts := make([]model.Part, 0, len(ids))

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, id := range ids {
		rec, ok := s.parts[id.String()]
		if !ok {
			return nil, fmt.Errorf("%w: %s", errs.ErrPartNotFound, id)
		}

		parts = append(parts, converter.PartFromRecord(rec))
	}

	return parts, nil
}
