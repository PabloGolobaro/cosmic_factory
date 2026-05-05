package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orderapi "github.com/PabloGolobaro/cosmic_factory/order/internal/api/order/v1"
	inventory "github.com/PabloGolobaro/cosmic_factory/order/internal/client/grpc/inventory/v1"
	payment "github.com/PabloGolobaro/cosmic_factory/order/internal/client/grpc/payment/v1"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/config"
	ordrepo "github.com/PabloGolobaro/cosmic_factory/order/internal/repository/order"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/orderitem"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/service/order"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

func main() {
	if err := run(); err != nil {
		slog.Error("ошибка запуска сервера", "error", err)
		os.Exit(1)
	}
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

func initPool(dsn string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("создание пула соединений: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("проверка соединения с БД: %w", err)
	}

	slog.Info("подключение к PostgreSQL установлено")

	return pool, nil
}

func run() error {
	if err := godotenv.Load("./../order.env"); err != nil {
		return fmt.Errorf("загрузка .env: %w", err)
	}

	configPath := config.ResolveConfigPath()

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	slog.Info("конфигурация загружена",
		"config_path", configPath,
		"http_address", cfg.HTTP.Address(),
		"grpc_address", cfg.GRPC.Address(),
		"pg_host", cfg.PG.Host,
		"inventory_address", cfg.Inventory.Address(),
		"payment_address", cfg.Payment.Address(),
	)

	inventoryConn, err := newGRPCConn(cfg.Inventory.Address(), cfg.Inventory.PingInterval, cfg.Inventory.PingTimeout)
	if err != nil {
		return fmt.Errorf("подключение к InventoryService: %w", err)
	}
	defer inventoryConn.Close()

	paymentConn, err := newGRPCConn(cfg.Payment.Address(), cfg.Payment.PingInterval, cfg.Payment.PingTimeout)
	if err != nil {
		return fmt.Errorf("подключение к PaymentService: %w", err)
	}
	defer paymentConn.Close()

	pool, err := initPool(cfg.PG.DSN())
	if err != nil {
		return err
	}
	defer pool.Close()

	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		return fmt.Errorf("создание transaction manager: %w", err)
	}

	orderRepo := ordrepo.NewOrderRepo(pool)
	orderItemRepo := orderitem.NewOrderItemRepo(pool)
	inventoryClient := inventory.NewInventoryClient(inventoryv1.NewInventoryServiceClient(inventoryConn))
	paymentClient := payment.NewPaymentClient(paymentv1.NewPaymentServiceClient(paymentConn))
	orderService := order.NewService(txManager, orderRepo, inventoryClient, paymentClient, orderItemRepo)
	orderApi := orderapi.NewApi(orderService)

	r, err := orderApi.SetupRouter()
	if err != nil {
		return fmt.Errorf("инициализация роутера: %w", err)
	}

	server := &http.Server{
		Addr:              cfg.HTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
	}

	go func() {
		slog.Info("запуск OrderService", "addr", cfg.HTTP.Address())
		if serveErr := server.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			slog.Error("ошибка запуска сервера", "error", serveErr)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	slog.Info("⚠️ Получен сигнал закрытия сервера. Выполняем graceful shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("остановка HTTP-сервера: %w", err)
	}

	slog.Info("✅ Сервер остановлен")

	return nil
}
