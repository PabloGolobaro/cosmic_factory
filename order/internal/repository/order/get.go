package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s *repo) Get(ctx context.Context, id uuid.UUID) (model.Order, error) {
	return s.fetchOrder(ctx, id, "SELECT * FROM orders WHERE uuid = $1;")
}
