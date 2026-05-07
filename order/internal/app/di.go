package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orderapi "github.com/PabloGolobaro/cosmic_factory/order/internal/api/order/v1"
	inventoryclient "github.com/PabloGolobaro/cosmic_factory/order/internal/client/grpc/inventory/v1"
	paymentclient "github.com/PabloGolobaro/cosmic_factory/order/internal/client/grpc/payment/v1"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/config"
	ordrepo "github.com/PabloGolobaro/cosmic_factory/order/internal/repository/order"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/orderitem"
	orderservice "github.com/PabloGolobaro/cosmic_factory/order/internal/service/order"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/closer"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

// diContainer — контейнер зависимостей (Composition Root) приложения.
//
// Каждый геттер следует паттерну ленивой инициализации (lazy initialization):
//  1. Проверяет, создан ли уже объект (nil-check).
//  2. Если нет — создаёт, запоминает в поле и возвращает.
//  3. Если да — сразу возвращает ранее созданный экземпляр.
type diContainer struct {
	conf config.Config

	// Инфраструктура (конкретные типы)
	pgPool        *pgxpool.Pool
	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn

	// Сервисный слой (интерфейсы из service/order/deps.go)
	txManager     orderservice.TxManager
	orderRepo     orderservice.OrderRepository
	orderItemRepo orderservice.OrderItemRepository
	invClient     orderservice.InventoryClient
	payClient     orderservice.PaymentClient

	// API-слой (интерфейс из api/order/v1/deps.go)
	orderSvc orderapi.OrderService

	// Презентационный слой
	router chi.Router
}

func newDIContainer(conf config.Config) *diContainer {
	return &diContainer{conf: conf}
}

// PGPool возвращает пул подключений к PostgreSQL.
func (d *diContainer) PGPool(ctx context.Context) (*pgxpool.Pool, error) {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, d.conf.PG.DSN())
		if err != nil {
			return nil, fmt.Errorf("создание пула соединений: %w", err)
		}

		if err = pool.Ping(ctx); err != nil {
			pool.Close()
			return nil, fmt.Errorf("ping PostgreSQL: %w", err)
		}

		slog.Info("подключение к PostgreSQL установлено")

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool, nil
}

// InventoryConn возвращает gRPC-соединение с сервисом Inventory.
func (d *diContainer) InventoryConn() (*grpc.ClientConn, error) {
	if d.inventoryConn == nil {
		conn, err := newGRPCConn(d.conf.Inventory.Address(), d.conf.Inventory.PingInterval, d.conf.Inventory.PingTimeout)
		if err != nil {
			return nil, fmt.Errorf("подключение к InventoryService: %w", err)
		}

		closer.Add("inventory gRPC connection", func(_ context.Context) error {
			return conn.Close()
		})

		d.inventoryConn = conn
	}

	return d.inventoryConn, nil
}

// PaymentConn возвращает gRPC-соединение с сервисом Payment.
func (d *diContainer) PaymentConn() (*grpc.ClientConn, error) {
	if d.paymentConn == nil {
		conn, err := newGRPCConn(d.conf.Payment.Address(), d.conf.Payment.PingInterval, d.conf.Payment.PingTimeout)
		if err != nil {
			return nil, fmt.Errorf("подключение к PaymentService: %w", err)
		}

		closer.Add("payment gRPC connection", func(_ context.Context) error {
			return conn.Close()
		})

		d.paymentConn = conn
	}

	return d.paymentConn, nil
}

// TxManager возвращает менеджер транзакций.
func (d *diContainer) TxManager(ctx context.Context) (orderservice.TxManager, error) {
	if d.txManager == nil {
		pool, err := d.PGPool(ctx)
		if err != nil {
			return nil, fmt.Errorf("transaction manager: %w", err)
		}

		txm, err := manager.New(trmpgx.NewDefaultFactory(pool))
		if err != nil {
			return nil, fmt.Errorf("создание transaction manager: %w", err)
		}

		d.txManager = txm
	}

	return d.txManager, nil
}

// OrderRepo возвращает репозиторий заказов.
func (d *diContainer) OrderRepo(ctx context.Context) (orderservice.OrderRepository, error) {
	if d.orderRepo == nil {
		pool, err := d.PGPool(ctx)
		if err != nil {
			return nil, fmt.Errorf("order repository: %w", err)
		}

		d.orderRepo = ordrepo.NewOrderRepo(pool)
	}

	return d.orderRepo, nil
}

// OrderItemRepo возвращает репозиторий позиций заказа.
func (d *diContainer) OrderItemRepo(ctx context.Context) (orderservice.OrderItemRepository, error) {
	if d.orderItemRepo == nil {
		pool, err := d.PGPool(ctx)
		if err != nil {
			return nil, fmt.Errorf("order item repository: %w", err)
		}

		d.orderItemRepo = orderitem.NewOrderItemRepo(pool)
	}

	return d.orderItemRepo, nil
}

// InvClient возвращает клиент сервиса Inventory.
func (d *diContainer) InvClient() (orderservice.InventoryClient, error) {
	if d.invClient == nil {
		conn, err := d.InventoryConn()
		if err != nil {
			return nil, fmt.Errorf("inventory client: %w", err)
		}

		d.invClient = inventoryclient.NewInventoryClient(inventoryv1.NewInventoryServiceClient(conn))
	}

	return d.invClient, nil
}

// PayClient возвращает клиент сервиса Payment.
func (d *diContainer) PayClient() (orderservice.PaymentClient, error) {
	if d.payClient == nil {
		conn, err := d.PaymentConn()
		if err != nil {
			return nil, fmt.Errorf("payment client: %w", err)
		}

		d.payClient = paymentclient.NewPaymentClient(paymentv1.NewPaymentServiceClient(conn))
	}

	return d.payClient, nil
}

// OrderService возвращает сервис бизнес-логики заказов.
func (d *diContainer) OrderService(ctx context.Context) (orderapi.OrderService, error) {
	if d.orderSvc == nil {
		txm, err := d.TxManager(ctx)
		if err != nil {
			return nil, fmt.Errorf("order service: %w", err)
		}

		orderRepo, err := d.OrderRepo(ctx)
		if err != nil {
			return nil, fmt.Errorf("order service: %w", err)
		}

		orderItemRepo, err := d.OrderItemRepo(ctx)
		if err != nil {
			return nil, fmt.Errorf("order service: %w", err)
		}

		invClient, err := d.InvClient()
		if err != nil {
			return nil, fmt.Errorf("order service: %w", err)
		}

		payClient, err := d.PayClient()
		if err != nil {
			return nil, fmt.Errorf("order service: %w", err)
		}

		d.orderSvc = orderservice.NewService(txm, orderRepo, invClient, payClient, orderItemRepo)
	}

	return d.orderSvc, nil
}

// Router возвращает настроенный HTTP-роутер приложения.
func (d *diContainer) Router(ctx context.Context) (chi.Router, error) {
	if d.router == nil {
		svc, err := d.OrderService(ctx)
		if err != nil {
			return nil, fmt.Errorf("router: %w", err)
		}

		r, err := orderapi.NewApi(svc).SetupRouter()
		if err != nil {
			return nil, fmt.Errorf("инициализация роутера: %w", err)
		}

		d.router = r
	}

	return d.router, nil
}

func newGRPCConn(addr string, pingInterval, pingTimeout time.Duration) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                pingInterval,
			Timeout:             pingTimeout,
			PermitWithoutStream: true,
		}))
}
