package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"buf.build/go/protovalidate"
	protovalidateMiddleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	v1 "github.com/PabloGolobaro/cosmic_factory/inventory/internal/api/part/v1"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/part"
	partSvc "github.com/PabloGolobaro/cosmic_factory/inventory/internal/service/part"
	"github.com/PabloGolobaro/cosmic_factory/shared/pkg/interceptors"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

const grpcAddress = ":50051"

const (
	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
)

func main() {
	if err := run(); err != nil {
		slog.Error("ошибка запуска сервера", "error", err)
		os.Exit(1)
	}
}

func run() error {
	lis, err := new(net.ListenConfig).Listen(context.Background(), "tcp", grpcAddress)
	if err != nil {
		return fmt.Errorf("создание listener: %w", err)
	}

	validator, err := protovalidate.New()
	if err != nil {
		return fmt.Errorf("создание валидатора: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: true,
		}),
		grpc.ChainUnaryInterceptor(
			interceptors.RecoveryInterceptor(),
			interceptors.LoggerInterceptor(),
			protovalidateMiddleware.UnaryServerInterceptor(validator),
		),
	)

	if err = godotenv.Load("./../../inventory.env"); err != nil {
		return fmt.Errorf("загрузка .env: %w", err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("DB_URI"))
	if err != nil {
		return fmt.Errorf("создание пула соединений: %w", err)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		return fmt.Errorf("проверка соединения с БД: %w", err)
	}
	slog.Info("подключение к PostgreSQL установлено")

	store := part.NewPartStore(pool)
	svc := partSvc.NewPartService(store)
	api := v1.NewApi(svc)

	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)
	reflection.Register(grpcServer)

	go func() {
		slog.Info("запуск InventoryService", "адрес", grpcAddress)
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			slog.Error("ошибка запуска сервера", "error", serveErr)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	slog.Info("⚠️ Получен сигнал закрытия сервера. Выполняем graceful shutdown")
	grpcServer.GracefulStop()
	slog.Info("✅ Сервер остановлен")

	return nil
}
