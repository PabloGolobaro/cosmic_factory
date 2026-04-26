package part

import (
	"context"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/converter"
)

func (s *store) GetAll(_ context.Context) ([]model.Part, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	parts := make([]model.Part, 0, len(s.parts))
	for _, rec := range s.parts {
		parts = append(parts, converter.PartFromRecord(rec))
	}
	return parts, nil
}
