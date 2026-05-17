package orderpaid

import (
	"context"
	"log/slog"
)

type service struct {
	consumer        Consumer
	assemblyService AssemblyService
}

func NewService(consumer Consumer, assemblyService AssemblyService) *service {
	return &service{
		consumer:        consumer,
		assemblyService: assemblyService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	slog.InfoContext(ctx, "запуск потребителя OrderPaid")

	return s.consumer.Consume(ctx, s.orderPaidHandler)
}
