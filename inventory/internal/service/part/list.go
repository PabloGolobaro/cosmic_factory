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
		parsedIDs := make([]uuid.UUID, 0, len(ids))
		for _, id := range ids {
			parsedID, err := uuid.Parse(id)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", errs.ErrInvalidUUID, err)
			}
			parsedIDs = append(parsedIDs, parsedID)
		}
		return s.PartRepository.GetBatch(ctx, parsedIDs)
	}

	parts, err := s.PartRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if partType != valueobject.PartTypeUnspecified {
		parts = slices.DeleteFunc(parts, func(p entity.Part) bool {
			return p.PartType() != partType
		})
	}

	slices.SortFunc(parts, func(a, b entity.Part) int {
		return cmp.Compare(strings.ToLower(a.Name()), strings.ToLower(b.Name()))
	})

	return parts, nil
}
