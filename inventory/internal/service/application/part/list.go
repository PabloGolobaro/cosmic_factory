package part

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
)

func (s service) List(ctx context.Context, ids []string, partType valueobject.PartType) ([]entity.Part, error) {
	if len(ids) > 0 {
		for _, id := range ids {
			if _, err := uuid.Parse(id); err != nil {
				return nil, fmt.Errorf("%w: %w", errs.ErrInvalidUUID, err)
			}
		}
		return s.PartRepository.GetBatch(ctx, valueobject.PartFilter{UUIDs: ids})
	}

	parts, err := s.PartRepository.GetBatch(ctx, valueobject.PartFilter{PartType: partType})
	if err != nil {
		return nil, err
	}

	slices.SortFunc(parts, func(a, b entity.Part) int {
		return cmp.Compare(strings.ToLower(a.Name()), strings.ToLower(b.Name()))
	})

	return parts, nil
}
