package assembly

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/assembly/internal/model"
)

func (s *service) Assemble(ctx context.Context, event model.OrderPaidEvent) error {
	buildSec := int64(5 + rand.IntN(11)) //nolint:mnd // диапазон [5, 15] секунд

	slog.InfoContext(ctx, "начало сборки корабля",
		"order_uuid", event.OrderUUID,
		"user_uuid", event.UserUUID,
		"build_sec", buildSec,
	)

	select {
	case <-time.After(time.Duration(buildSec) * time.Second):
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
