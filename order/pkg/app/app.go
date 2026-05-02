package app

import (
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	orderapi "github.com/PabloGolobaro/cosmic_factory/order/internal/api/order/v1"
	inventory "github.com/PabloGolobaro/cosmic_factory/order/internal/client/grpc/inventory/v1"
	payment "github.com/PabloGolobaro/cosmic_factory/order/internal/client/grpc/payment/v1"
	ordrepo "github.com/PabloGolobaro/cosmic_factory/order/internal/repository/order"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/orderitem"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/service/order"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

func NewHTTPHandler(pool *pgxpool.Pool, txManager *manager.Manager, inventoryServiceClient inventoryv1.InventoryServiceClient, paymentServiceClient paymentv1.PaymentServiceClient) (chi.Router, error) {
	orderRepo := ordrepo.NewOrderRepo(pool)
	orderItemRepo := orderitem.NewOrderItemRepo(pool)

	inventoryClient := inventory.NewInventoryClient(inventoryServiceClient)
	paymentClient := payment.NewPaymentClient(paymentServiceClient)

	orderService := order.NewService(txManager, orderRepo, inventoryClient, paymentClient, orderItemRepo)

	orderApi := orderapi.NewApi(orderService)

	r, err := orderApi.SetupRouter()
	if err != nil {
		slog.Error("Не удалось инициализировать роутер", "error", err)
	}

	return r, err
}
