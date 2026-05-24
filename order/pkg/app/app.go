package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	orderapi "github.com/PabloGolobaro/cosmic_factory/order/internal/api/order/v1"
	iamv1 "github.com/PabloGolobaro/cosmic_factory/order/internal/client/grpc/iam/v1"
	inventory "github.com/PabloGolobaro/cosmic_factory/order/internal/client/grpc/inventory/v1"
	payment "github.com/PabloGolobaro/cosmic_factory/order/internal/client/grpc/payment/v1"
	authmw "github.com/PabloGolobaro/cosmic_factory/order/internal/middleware"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	ordrepo "github.com/PabloGolobaro/cosmic_factory/order/internal/repository/order"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/orderitem"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/service/order"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/auth"
	authv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/auth/v1"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

// NewHTTPHandler создаёт HTTP-роутер с noop-продюсером и auth middleware (для API-тестов).
func NewHTTPHandler(pool *pgxpool.Pool, txManager *manager.Manager, inventoryServiceClient inventoryv1.InventoryServiceClient, paymentServiceClient paymentv1.PaymentServiceClient, authServiceClient authv1.AuthServiceClient) (chi.Router, error) {
	return buildRouter(pool, txManager, inventoryServiceClient, paymentServiceClient, authmw.New(iamv1.New(authServiceClient)), noopProducer{})
}

// NewHTTPHandlerWithProducer создаёт HTTP-роутер с реальным Kafka-продюсером (для e2e-тестов).
// Использует заглушку auth — в e2e тестируется Kafka-цепочка, а не аутентификация.
func NewHTTPHandlerWithProducer(pool *pgxpool.Pool, txManager *manager.Manager, inventoryServiceClient inventoryv1.InventoryServiceClient, paymentServiceClient paymentv1.PaymentServiceClient, orderPaidProducer order.OrderPaidProducer) (chi.Router, error) {
	return buildRouter(pool, txManager, inventoryServiceClient, paymentServiceClient, e2eAuthMiddleware, orderPaidProducer)
}

// e2eAuthMiddleware инжектирует случайный user UUID в контекст без вызова IAM.
// Предназначен исключительно для e2e-тестов, где предмет проверки — Kafka-цепочка.
func e2eAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := auth.WithUserUUID(r.Context(), uuid.New())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func buildRouter(pool *pgxpool.Pool, txManager *manager.Manager, inventoryServiceClient inventoryv1.InventoryServiceClient, paymentServiceClient paymentv1.PaymentServiceClient, authMiddleware func(http.Handler) http.Handler, orderPaidProducer order.OrderPaidProducer) (chi.Router, error) {
	orderRepo := ordrepo.NewOrderRepo(pool)
	orderItemRepo := orderitem.NewOrderItemRepo(pool)

	inventoryClient := inventory.NewInventoryClient(inventoryServiceClient)
	paymentClient := payment.NewPaymentClient(paymentServiceClient)

	orderService := order.NewService(txManager, orderRepo, inventoryClient, paymentClient, orderItemRepo, orderPaidProducer)

	orderApi := orderapi.NewApi(orderService)

	r, err := orderApi.SetupRouter(authMiddleware)
	if err != nil {
		slog.Error("Не удалось инициализировать роутер", "error", err)
	}

	return r, err
}

type noopProducer struct{}

func (noopProducer) PublishOrderPaid(_ context.Context, _ model.OrderPaidEvent) error { return nil }
