package orderpaid

import (
	"context"
	"log/slog"

	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka"
)

func (s *service) orderPaidHandler(ctx context.Context, msg kafka.Message) error {
	event, err := decodeOrderPaid(msg.Value)
	if err != nil {
		slog.ErrorContext(ctx, "не удалось декодировать OrderPaid", "error", err)
		return err
	}

	slog.InfoContext(ctx, "получено событие OrderPaid",
		"topic", msg.Topic,
		"partition", msg.Partition,
		"offset", msg.Offset,
		"order_uuid", event.OrderUUID,
		"user_uuid", event.UserUUID,
	)

	return s.assemblyService.Assemble(ctx, event)
}
