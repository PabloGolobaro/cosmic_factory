package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/converter"
)

func (s *store) Get(_ context.Context, id uuid.UUID) (model.Part, error) {
	s.mu.RLock()
	rec, ok := s.parts[id.String()]
	s.mu.RUnlock()

	if !ok {
		return model.Part{}, fmt.Errorf("%w: %s", errs.ErrPartNotFound, id)
	}

	return converter.PartFromRecord(rec), nil
}
