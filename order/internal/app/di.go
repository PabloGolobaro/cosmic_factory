package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
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
	shipassembled "github.com/PabloGolobaro/cosmic_factory/order/internal/consumer/ship_assembled"
	orderpaidproducer "github.com/PabloGolobaro/cosmic_factory/order/internal/producer/order"
	ordrepo "github.com/PabloGolobaro/cosmic_factory/order/internal/repository/order"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/orderitem"
	orderservice "github.com/PabloGolobaro/cosmic_factory/order/internal/service/order"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/closer"
	kafkaconsumer "github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka/consumer"
	kafkaproducer "github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka/producer"
	kafkamw "github.com/PabloGolobaro/cosmic_factory/platform/pkg/middleware/kafka"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

type consumerRunner interface {
	RunConsumer(ctx context.Context) error
}

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
	consumerGroup sarama.ConsumerGroup
	syncProducer  sarama.SyncProducer

	// Сервисный слой (интерфейсы из service/order/deps.go)
	txManager        orderservice.TxManager
	orderRepo        orderservice.OrderRepository
	orderItemRepo    orderservice.OrderItemRepository
	invClient        orderservice.InventoryClient
	payClient        orderservice.PaymentClient
	orderPaidProd    orderservice.OrderPaidProducer
	shipAssembledSvc shipassembled.ShipAssembledService

	// API-слой (интерфейс из api/order/v1/deps.go)
	orderSvc orderapi.OrderService

	// Kafka consumer runner
	shipAssembledRunner consumerRunner

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

// KafkaConsumerGroup возвращает Kafka consumer group.
func (d *diContainer) KafkaConsumerGroup() (sarama.ConsumerGroup, error) {
	if d.consumerGroup == nil {
		cfg := sarama.NewConfig()
		cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
		cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
			sarama.NewBalanceStrategyRoundRobin(),
		}

		group, err := sarama.NewConsumerGroup(d.conf.Kafka.Brokers, d.conf.Kafka.ConsumerGroup, cfg)
		if err != nil {
			return nil, fmt.Errorf("создание Kafka consumer group: %w", err)
		}

		slog.Info("Kafka consumer group создана", "group", d.conf.Kafka.ConsumerGroup)

		closer.Add("kafka consumer group", func(_ context.Context) error {
			return group.Close()
		})

		d.consumerGroup = group
	}

	return d.consumerGroup, nil
}

// KafkaSyncProducer возвращает Kafka sync producer.
func (d *diContainer) KafkaSyncProducer() (sarama.SyncProducer, error) {
	if d.syncProducer == nil {
		cfg := sarama.NewConfig()
		cfg.Producer.Return.Successes = true
		cfg.Producer.RequiredAcks = sarama.WaitForAll

		producer, err := sarama.NewSyncProducer(d.conf.Kafka.Brokers, cfg)
		if err != nil {
			return nil, fmt.Errorf("создание Kafka sync producer: %w", err)
		}

		slog.Info("Kafka sync producer создан", "topic", d.conf.Kafka.ProduceTopic)

		closer.Add("kafka sync producer", func(_ context.Context) error {
			return producer.Close()
		})

		d.syncProducer = producer
	}

	return d.syncProducer, nil
}

// OrderPaidProducer возвращает продюсер события OrderPaid.
func (d *diContainer) OrderPaidProducer() (orderservice.OrderPaidProducer, error) {
	if d.orderPaidProd == nil {
		syncProducer, err := d.KafkaSyncProducer()
		if err != nil {
			return nil, fmt.Errorf("order paid producer: %w", err)
		}

		p := kafkaproducer.NewProducer(syncProducer, d.conf.Kafka.ProduceTopic)
		d.orderPaidProd = orderpaidproducer.NewService(p)
	}

	return d.orderPaidProd, nil
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

		orderPaidProd, err := d.OrderPaidProducer()
		if err != nil {
			return nil, fmt.Errorf("order service: %w", err)
		}

		svc := orderservice.NewService(txm, orderRepo, invClient, payClient, orderItemRepo, orderPaidProd)
		d.orderSvc = svc
		d.shipAssembledSvc = svc
	}

	return d.orderSvc, nil
}

// ShipAssembledConsumerService возвращает consumer ShipAssembled событий.
func (d *diContainer) ShipAssembledConsumerService() (consumerRunner, error) {
	if d.shipAssembledRunner == nil {
		group, err := d.KafkaConsumerGroup()
		if err != nil {
			return nil, fmt.Errorf("ship assembled consumer: %w", err)
		}

		if d.shipAssembledSvc == nil {
			_, err = d.OrderService(context.Background())
			if err != nil {
				return nil, fmt.Errorf("ship assembled consumer: %w", err)
			}
		}

		consumer := kafkaconsumer.NewConsumer(
			group,
			[]string{d.conf.Kafka.ConsumeTopic},
			kafkaconsumer.WithMiddlewares(kafkamw.ConsumerLogging()),
		)

		d.shipAssembledRunner = shipassembled.NewService(consumer, d.shipAssembledSvc)
	}

	return d.shipAssembledRunner, nil
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
