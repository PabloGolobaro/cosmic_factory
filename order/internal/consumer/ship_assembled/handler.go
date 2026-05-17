package shipassembled

import (
	"context"
	"log/slog"

	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka"
)

func (s *service) shipAssembledHandler(ctx context.Context, msg kafka.Message) error {
	event, err := decodeShipAssembled(msg.Value)
	if err != nil {
		slog.ErrorContext(ctx, "не удалось декодировать ShipAssembled", "error", err)
		return err
	}

	slog.InfoContext(ctx, "получено событие ShipAssembled",
		"topic", msg.Topic,
		"partition", msg.Partition,
		"offset", msg.Offset,
		"order_uuid", event.OrderUUID,
		"build_time_sec", event.BuildTimeSec,
	)

	return s.shipAssembledService.CommitShipParts(ctx, event.OrderUUID)
}
