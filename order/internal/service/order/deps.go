package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

type Repository interface {
	Create(context.Context, model.Order) (model.Order, error)
	Get(context.Context, uuid.UUID) (model.Order, error)
	Delete(context.Context, uuid.UUID) error
	Update(context.Context, model.Order) error
}

type InventoryClient interface {
	ListParts(ctx context.Context, uuids []string) ([]model.Part, error)
}
type PaymentClient interface {
	PayOrder(ctx context.Context, uuid string, paymentMethod model.PaymentMethod) error
}
