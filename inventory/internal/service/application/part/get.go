package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
)

func (s service) Get(ctx context.Context, id string) (entity.Part, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return entity.Part{}, fmt.Errorf("%w: %w", errs.ErrInvalidUUID, err)
	}
	return s.PartRepository.Get(ctx, parsedID)
}
