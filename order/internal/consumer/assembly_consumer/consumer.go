package assemblyconsumer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka"
)

type service struct {
	consumer        Consumer
	orderRepo       OrderRepository
	inventoryClient InventoryClient
	txManager       TxManager
}

func NewService(consumer Consumer, orderRepo OrderRepository, invClient InventoryClient, txManager TxManager) *service {
	return &service{
		consumer:        consumer,
		orderRepo:       orderRepo,
		inventoryClient: invClient,
		txManager:       txManager,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	slog.InfoContext(ctx, "запуск потребителя ShipAssembled")

	return s.consumer.Consume(ctx, s.shipAssembledHandler)
}

func (s *service) shipAssembledHandler(ctx context.Context, msg kafka.Message) error {
	event, err := decodeShipAssembled(msg.Value)
	if err != nil {
		slog.ErrorContext(ctx, "не удалось декодировать ShipAssembled", "error", err)
		return err
	}

	slog.InfoContext(ctx, "получено событие ShipAssembled",
		"order_uuid", event.OrderUUID,
		"build_time_sec", event.BuildTimeSec,
	)

	return s.commitShipParts(ctx, event.OrderUUID)
}

func (s *service) commitShipParts(ctx context.Context, orderUUID string) error {
	id, err := uuid.Parse(orderUUID)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrInvalidUUID, err)
	}

	return s.txManager.Do(ctx, func(txCtx context.Context) error {
		order, err := s.orderRepo.GetForUpdate(txCtx, id)
		if err != nil {
			return fmt.Errorf("%w: %w", errs.ErrOrderNotFound, err)
		}

		if order.Status == model.OrderStatusAssembled {
			return nil
		}

		uuids := []string{order.HullUUID.String(), order.EngineUUID.String()}
		if order.ShieldUUID != nil {
			uuids = append(uuids, order.ShieldUUID.String())
		}
		if order.WeaponUUID != nil {
			uuids = append(uuids, order.WeaponUUID.String())
		}

		if err = s.inventoryClient.CommitParts(txCtx, uuids); err != nil {
			return err
		}

		order.Status = model.OrderStatusAssembled
		return s.orderRepo.Update(txCtx, order)
	})
}
