package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/converter"
)

func (s *store) Get(_ context.Context, id uuid.UUID) (model.Order, error) {
	s.mu.RLock()
	rec, ok := s.orders[id.String()]
	s.mu.RUnlock()

	if !ok {
		return model.Order{}, fmt.Errorf("%w: %s", errs.ErrOrderNotFound, id)
	}

	return converter.OrderFromRecord(rec), nil
}
