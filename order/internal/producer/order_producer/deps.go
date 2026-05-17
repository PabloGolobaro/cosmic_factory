package orderproducer

import (
	"context"

	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka"
)

// KafkaProducer определяет контракт для отправки сообщений в Kafka.
type KafkaProducer interface {
	Send(ctx context.Context, msg *kafka.Message) error
}
