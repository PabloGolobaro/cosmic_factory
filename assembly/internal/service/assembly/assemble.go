package assembly

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/assembly/internal/model"
)

func (s *service) Assemble(ctx context.Context, event model.OrderPaidEvent) error {
	delay := s.getBuildDelay()
	buildSec := int64(delay.Seconds())

	slog.InfoContext(ctx, "начало сборки корабля",
		"order_uuid", event.OrderUUID,
		"user_uuid", event.UserUUID,
		"build_sec", buildSec,
	)

	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return ctx.Err()
	}

	assembled := model.ShipAssembledEvent{
		EventUUID:    uuid.New().String(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: buildSec,
		AssembledAt:  time.Now(),
	}

	slog.InfoContext(ctx, "сборка завершена",
		"order_uuid", event.OrderUUID,
		"build_sec", buildSec,
	)

	return s.producer.PublishShipAssembled(ctx, assembled)
}
