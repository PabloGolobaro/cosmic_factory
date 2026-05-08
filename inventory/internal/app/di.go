package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	apipart "github.com/PabloGolobaro/cosmic_factory/inventory/internal/api/part/v1"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/config"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/part"
	part2 "github.com/PabloGolobaro/cosmic_factory/inventory/internal/service/application/part"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/service/domain"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/closer"
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
	pgPool *pgxpool.Pool

	// Репозиторный слой (интерфейс из service/part/deps.go)
	partRepo part2.PartRepository

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

// PartSvc возвращает сервис бизнес-логики деталей.
func (d *diContainer) PartSvc(ctx context.Context) (apipart.PartService, error) {
	if d.partSvc == nil {
		repo, err := d.PartRepo(ctx)
		if err != nil {
			return nil, fmt.Errorf("part service: %w", err)
		}

		d.partSvc = part2.NewPartService(repo, domain.NewCompatibilityChecker())
	}

	return d.partSvc, nil
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
