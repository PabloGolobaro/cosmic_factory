package order

import (
	"context"
	"fmt"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/converter"
)

func (s *store) Update(_ context.Context, order model.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := order.OrderUUID.String()
	if _, ok := s.orders[key]; !ok {
		return fmt.Errorf("%w: %s", errs.ErrOrderNotFound, order.OrderUUID)
	}

	s.orders[key] = converter.OrderToRecord(order)
	return nil
}
