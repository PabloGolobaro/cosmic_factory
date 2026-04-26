package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/converter"
)

func (s *store) Create(_ context.Context, order model.Order) (model.Order, error) {
	order.OrderUUID = uuid.New()

	rec := converter.OrderToRecord(order)

	s.mu.Lock()
	s.orders[order.OrderUUID.String()] = rec
	s.mu.Unlock()

	return order, nil
}
