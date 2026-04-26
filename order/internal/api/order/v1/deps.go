package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

type OrderService interface {
	Create(context.Context, model.Order) (model.Order, error)
	Get(context.Context, string) (*model.Order, error)
	Cancel(context.Context, string) error
	Pay(context.Context, string, string) (uuid.UUID, error)
}
