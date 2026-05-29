package order

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type service struct {
	txManager           TxManager
	Repository          OrderRepository
	InventoryClient     InventoryClient
	PaymentClient       PaymentClient
	OrderItemRepository OrderItemRepository
	OrderPaidProducer   OrderPaidProducer

	ordersCreated metric.Int64Counter
	ordersPaid    metric.Int64Counter
	ordersRevenue metric.Int64Counter
}

func NewService(txManager TxManager, repository OrderRepository, inventoryClient InventoryClient, paymentClient PaymentClient, orderItemRepository OrderItemRepository, orderPaidProducer OrderPaidProducer) (*service, error) {
	meter := otel.Meter("orders")

	ordersCreated, err := meter.Int64Counter("orders_created",
		metric.WithDescription("Количество созданных заказов"))
	if err != nil {
		return nil, fmt.Errorf("orders_created counter: %w", err)
	}
	ordersPaid, err := meter.Int64Counter("orders_paid",
		metric.WithDescription("Количество оплаченных заказов"))
	if err != nil {
		return nil, fmt.Errorf("orders_paid counter: %w", err)
	}
	ordersRevenue, err := meter.Int64Counter("orders_revenue",
		metric.WithDescription("Суммарная выручка в копейках"))
	if err != nil {
		return nil, fmt.Errorf("orders_revenue counter: %w", err)
	}

	return &service{
		txManager:           txManager,
		Repository:          repository,
		InventoryClient:     inventoryClient,
		PaymentClient:       paymentClient,
		OrderItemRepository: orderItemRepository,
		OrderPaidProducer:   orderPaidProducer,
		ordersCreated:       ordersCreated,
		ordersPaid:          ordersPaid,
		ordersRevenue:       ordersRevenue,
	}, nil
}
