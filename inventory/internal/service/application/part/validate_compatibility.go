package part

import (
	"context"
	"fmt"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
)

func (s *service) ValidateCompatibility(ctx context.Context, uuids []string) error {
	parts, err := s.PartRepository.GetBatch(ctx, valueobject.PartFilter{UUIDs: uuids})
	if err != nil {
		return fmt.Errorf("получить детали: %w", err)
	}

	return s.compatibilityChecker.Check(parts)
}
