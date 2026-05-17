package order

import (
	"context"

	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka"
)

// Consumer определяет контракт для потребления сообщений из Kafka
type Consumer interface {
	Consume(ctx context.Context, handler kafka.MessageHandler) error
}
