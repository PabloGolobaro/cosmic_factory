package shipassembled

import (
	"context"

	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka"
)

type KafkaProducer interface {
	Send(ctx context.Context, msg *kafka.Message) error
}
