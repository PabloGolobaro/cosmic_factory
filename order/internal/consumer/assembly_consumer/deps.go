package assemblyconsumer

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka"
)

type Consumer interface {
	Consume(ctx context.Context, handler kafka.MessageHandler) error
}

type OrderRepository interface {
	GetForUpdate(ctx context.Context, id uuid.UUID) (model.Order, error)
	Update(ctx context.Context, order model.Order) error
}

type InventoryClient interface {
	CommitParts(ctx context.Context, partUUIDs []string) error
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
