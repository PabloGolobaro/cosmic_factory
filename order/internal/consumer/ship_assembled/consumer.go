package shipassembled

import (
	"context"
	"log/slog"
)

type service struct {
	consumer             Consumer
	shipAssembledService ShipAssembledService
}

func NewService(consumer Consumer, svc ShipAssembledService) *service {
	return &service{
		consumer:             consumer,
		shipAssembledService: svc,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	slog.InfoContext(ctx, "запуск потребителя ShipAssembled")

	return s.consumer.Consume(ctx, s.shipAssembledHandler)
}
