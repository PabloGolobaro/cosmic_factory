package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	apiauth "github.com/PabloGolobaro/cosmic_factory/iam/internal/api/auth/v1"
	apiuser "github.com/PabloGolobaro/cosmic_factory/iam/internal/api/user/v1"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/config"
	reposes "github.com/PabloGolobaro/cosmic_factory/iam/internal/repository/session"
	repouser "github.com/PabloGolobaro/cosmic_factory/iam/internal/repository/user"
	svcauth "github.com/PabloGolobaro/cosmic_factory/iam/internal/service/auth"
	svcuser "github.com/PabloGolobaro/cosmic_factory/iam/internal/service/user"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/closer"
	platformredis "github.com/PabloGolobaro/cosmic_factory/platform/pkg/redis"
	authproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/auth/v1"
	userproto "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/user/v1"
)

// diContainer — контейнер зависимостей (Composition Root) приложения.
//
// Каждый геттер следует паттерну ленивой инициализации (lazy initialization).
type diContainer struct {
	conf config.Config

	// Инфраструктура
	pgPool      *pgxpool.Pool
	redisClient *redis.Client

	// Репозиторный слой
	userRepo    repouser.Repository
	sessionRepo reposes.Repository

	// Сервисный слой
	authSvc svcauth.Service
	userSvc svcuser.Service

	// gRPC-обработчики
	authHandler authproto.AuthServiceServer
	userHandler userproto.UserServiceServer
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

// RedisClient возвращает Redis-клиент.
func (d *diContainer) RedisClient(_ context.Context) (*redis.Client, error) {
	if d.redisClient == nil {
		client, err := platformredis.NewClient(
			&redis.Options{Addr: d.conf.Redis.Addr()},
			slog.Default(),
		)
		if err != nil {
			return nil, fmt.Errorf("создание Redis-клиента: %w", err)
		}

		closer.Add("Redis client", func(_ context.Context) error {
			return client.Close()
		})

		d.redisClient = client
	}

	return d.redisClient, nil
}

// UserRepo возвращает репозиторий пользователей.
func (d *diContainer) UserRepo(_ context.Context) (repouser.Repository, error) {
	if d.userRepo == nil {
		d.userRepo = repouser.NewStub()
	}

	return d.userRepo, nil
}

// SessionRepo возвращает репозиторий сессий.
func (d *diContainer) SessionRepo(_ context.Context) (reposes.Repository, error) {
	if d.sessionRepo == nil {
		d.sessionRepo = reposes.NewStub()
	}

	return d.sessionRepo, nil
}

// AuthSvc возвращает сервис аутентификации.
func (d *diContainer) AuthSvc(_ context.Context) (svcauth.Service, error) {
	if d.authSvc == nil {
		d.authSvc = svcauth.NewStub()
	}

	return d.authSvc, nil
}

// UserSvc возвращает сервис управления пользователями.
func (d *diContainer) UserSvc(_ context.Context) (svcuser.Service, error) {
	if d.userSvc == nil {
		d.userSvc = svcuser.NewStub()
	}

	return d.userSvc, nil
}

// AuthHandler возвращает gRPC-обработчик AuthService.
func (d *diContainer) AuthHandler(ctx context.Context) (authproto.AuthServiceServer, error) {
	if d.authHandler == nil {
		svc, err := d.AuthSvc(ctx)
		if err != nil {
			return nil, fmt.Errorf("auth handler: %w", err)
		}

		d.authHandler = apiauth.NewAPI(svc)
	}

	return d.authHandler, nil
}

// UserHandler возвращает gRPC-обработчик UserService.
func (d *diContainer) UserHandler(ctx context.Context) (userproto.UserServiceServer, error) {
	if d.userHandler == nil {
		svc, err := d.UserSvc(ctx)
		if err != nil {
			return nil, fmt.Errorf("user handler: %w", err)
		}

		d.userHandler = apiuser.NewAPI(svc)
	}

	return d.userHandler, nil
}
