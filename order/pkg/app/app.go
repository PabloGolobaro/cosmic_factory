package app

import (
	"context"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	orderapi "github.com/PabloGolobaro/cosmic_factory/order/internal/api/order/v1"
	inventory "github.com/PabloGolobaro/cosmic_factory/order/internal/client/grpc/inventory/v1"
	payment "github.com/PabloGolobaro/cosmic_factory/order/internal/client/grpc/payment/v1"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	ordrepo "github.com/PabloGolobaro/cosmic_factory/order/internal/repository/order"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/orderitem"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/service/order"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

// NewHTTPHandler создаёт HTTP-роутер с noop-продюсером (для интеграционных тестов без Kafka).
func NewHTTPHandler(pool *pgxpool.Pool, txManager *manager.Manager, inventoryServiceClient inventoryv1.InventoryServiceClient, paymentServiceClient paymentv1.PaymentServiceClient) (chi.Router, error) {
	return NewHTTPHandlerWithProducer(pool, txManager, inventoryServiceClient, paymentServiceClient, noopProducer{})
}

// NewHTTPHandlerWithProducer создаёт HTTP-роутер с реальным Kafka-продюсером (для e2e-тестов).
func NewHTTPHandlerWithProducer(pool *pgxpool.Pool, txManager *manager.Manager, inventoryServiceClient inventoryv1.InventoryServiceClient, paymentServiceClient paymentv1.PaymentServiceClient, orderPaidProducer order.OrderPaidProducer) (chi.Router, error) {
	orderRepo := ordrepo.NewOrderRepo(pool)
	orderItemRepo := orderitem.NewOrderItemRepo(pool)

	inventoryClient := inventory.NewInventoryClient(inventoryServiceClient)
	paymentClient := payment.NewPaymentClient(paymentServiceClient)

	orderService := order.NewService(txManager, orderRepo, inventoryClient, paymentClient, orderItemRepo, orderPaidProducer)

	orderApi := orderapi.NewApi(orderService)

	r, err := orderApi.SetupRouter()
	if err != nil {
		slog.Error("Не удалось инициализировать роутер", "error", err)
	}

	return r, err
}

type noopProducer struct{}

func (noopProducer) PublishOrderPaid(_ context.Context, _ model.OrderPaidEvent) error { return nil }
