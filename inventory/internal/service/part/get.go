package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
)

func (s service) Get(ctx context.Context, id string) (model.Part, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return model.Part{}, fmt.Errorf("%w: %w", errs.ErrInvalidUUID, err)
	}
	return s.PartRepository.Get(ctx, parsedID)
}
