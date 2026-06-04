package app

import (
	"context"
	"fmt"
	"log/slog"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	apipart "github.com/PabloGolobaro/cosmic_factory/inventory/internal/api/part/v1"
	iamv1client "github.com/PabloGolobaro/cosmic_factory/inventory/internal/client/grpc/iam/v1"
	parttracing "github.com/PabloGolobaro/cosmic_factory/inventory/internal/service/application/part/tracing"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/config"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/part"
	part2 "github.com/PabloGolobaro/cosmic_factory/inventory/internal/service/application/part"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/service/domain"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/closer"
	authproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/auth/v1"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

// diContainer — контейнер зависимостей (Composition Root) приложения.
//
// Каждый геттер следует паттерну ленивой инициализации (lazy initialization):
//  1. Проверяет, создан ли уже объект (nil-check).
//  2. Если нет — создаёт, запоминает в поле и возвращает.
//  3. Если да — сразу возвращает ранее созданный экземпляр.
type diContainer struct {
	conf config.Config

	// Инфраструктура
	pgPool  *pgxpool.Pool
	iamConn *grpc.ClientConn

	iamClient *iamv1client.Client

	// Репозиторный слой (интерфейс из service/part/deps.go)
	partRepo part2.PartRepository

	// Transaction manager
	txManager part2.TxManager

	// Сервисный слой (интерфейс из api/part/v1/deps.go)
	partSvc apipart.PartService

	// gRPC handler
	inventoryHandler inventoryv1.InventoryServiceServer
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

// PartRepo возвращает репозиторий деталей.
func (d *diContainer) PartRepo(ctx context.Context) (part2.PartRepository, error) {
	if d.partRepo == nil {
		pool, err := d.PGPool(ctx)
		if err != nil {
			return nil, fmt.Errorf("part repository: %w", err)
		}

		d.partRepo = part.NewPartStore(pool)
	}

	return d.partRepo, nil
}

// TxMgr возвращает менеджер транзакций.
func (d *diContainer) TxMgr(ctx context.Context) (part2.TxManager, error) {
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

// PartSvc возвращает сервис бизнес-логики деталей.
func (d *diContainer) PartSvc(ctx context.Context) (apipart.PartService, error) {
	if d.partSvc == nil {
		repo, err := d.PartRepo(ctx)
		if err != nil {
			return nil, fmt.Errorf("part service: %w", err)
		}

		txm, err := d.TxMgr(ctx)
		if err != nil {
			return nil, fmt.Errorf("part service: %w", err)
		}

		d.partSvc = parttracing.NewTracedService(part2.NewPartService(repo, domain.NewCompatibilityChecker(), txm))
	}

	return d.partSvc, nil
}

// IAMConn возвращает gRPC-соединение с сервисом IAM.
func (d *diContainer) IAMConn(_ context.Context) (*grpc.ClientConn, error) {
	if d.iamConn == nil {
		conn, err := grpc.NewClient(d.conf.IAM.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                d.conf.IAM.PingInterval,
				Timeout:             d.conf.IAM.PingTimeout,
				PermitWithoutStream: true,
			}),
		)
		if err != nil {
			return nil, fmt.Errorf("подключение к IAMService: %w", err)
		}

		closer.Add("iam gRPC connection", func(_ context.Context) error {
			return conn.Close()
		})

		d.iamConn = conn
	}

	return d.iamConn, nil
}

// IAMClient возвращает клиент сервиса IAM.
func (d *diContainer) IAMClient(ctx context.Context) (*iamv1client.Client, error) {
	if d.iamClient == nil {
		conn, err := d.IAMConn(ctx)
		if err != nil {
			return nil, fmt.Errorf("iam client: %w", err)
		}

		d.iamClient = iamv1client.New(authproto.NewAuthServiceClient(conn))
	}

	return d.iamClient, nil
}

// InventoryHandler возвращает gRPC-обработчик сервиса инвентаря.
func (d *diContainer) InventoryHandler(ctx context.Context) (inventoryv1.InventoryServiceServer, error) {
	if d.inventoryHandler == nil {
		svc, err := d.PartSvc(ctx)
		if err != nil {
			return nil, fmt.Errorf("inventory handler: %w", err)
		}

		d.inventoryHandler = apipart.NewApi(svc)
	}

	return d.inventoryHandler, nil
}
