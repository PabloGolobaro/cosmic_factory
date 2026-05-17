package shipassembled

import (
	"context"

	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka"
)

type Consumer interface {
	Consume(ctx context.Context, handler kafka.MessageHandler) error
}

type ShipAssembledService interface {
	CommitShipParts(ctx context.Context, orderUUID string) error
}
