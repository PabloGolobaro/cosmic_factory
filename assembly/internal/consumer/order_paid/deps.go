package orderpaid

import (
	"context"

	"github.com/PabloGolobaro/cosmic_factory/assembly/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka"
)

type Consumer interface {
	Consume(ctx context.Context, handler kafka.MessageHandler) error
}

type AssemblyService interface {
	Assemble(ctx context.Context, event model.OrderPaidEvent) error
}
