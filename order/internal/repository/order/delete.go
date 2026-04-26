package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
)

func (s *store) Delete(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := id.String()
	if _, ok := s.orders[key]; !ok {
		return fmt.Errorf("%w: %s", errs.ErrOrderNotFound, id)
	}

	delete(s.orders, key)
	return nil
}
