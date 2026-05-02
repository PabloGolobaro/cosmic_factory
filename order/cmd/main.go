package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
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
	ordrepo "github.com/PabloGolobaro/cosmic_factory/order/internal/repository/order"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/orderitem"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/service/order"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

const (
	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"
)

const (
	clientPingInterval = 10 * time.Second
	clientPingTimeout  = 10 * time.Second
)

const (
	httpPort = "8080"

	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("ошибка запуска сервера", "error", err)
		os.Exit(1)
	}
}

func newGRPCConn(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                clientPingInterval,
			Timeout:             clientPingTimeout,
			PermitWithoutStream: true,
		}))
}

func run() error {
	inventoryConn, err := newGRPCConn(inventoryServiceAddress)
	if err != nil {
		return fmt.Errorf("подключение к InventoryService: %w", err)
	}
	defer inventoryConn.Close()

	paymentConn, err := newGRPCConn(paymentServiceAddress)
	if err != nil {
		return fmt.Errorf("подключение к PaymentService: %w", err)
	}
	defer paymentConn.Close()

	if err = godotenv.Load("./../../order.env"); err != nil {
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
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	go func() {
		slog.Info("запуск OrderService", "port", httpPort)
		if serveErr := server.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			slog.Error("ошибка запуска сервера", "error", serveErr)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	slog.Info("⚠️ Получен сигнал закрытия сервера. Выполняем graceful shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("остановка HTTP-сервера: %w", err)
	}

	slog.Info("✅ Сервер остановлен")

	return nil
}
