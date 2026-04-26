package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-faster/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orderapi "github.com/PabloGolobaro/cosmic_factory/order/internal/api/order/v1"
	inventory "github.com/PabloGolobaro/cosmic_factory/order/internal/client/grpc/inventory/v1"
	payment "github.com/PabloGolobaro/cosmic_factory/order/internal/client/grpc/payment/v1"
	ordstore "github.com/PabloGolobaro/cosmic_factory/order/internal/repository/order"
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

	// Таймауты для HTTP-сервера.
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	// Создать gRPC соединение с InventoryService
	inventoryConn, err := grpc.NewClient(inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                clientPingInterval, // Интервал ping'ов для обнаружения мёртвых соединений
			Timeout:             clientPingTimeout,  // Таймаут ожидания pong
			PermitWithoutStream: true,               // Держать соединение "тёплым" без активных RPC
		}))
	if err != nil {
		slog.Error("не удалось подключиться к InventoryService", "error", err)
		return
	}
	defer inventoryConn.Close()

	// Создать gRPC соединение с PaymentService
	paymentConn, err := grpc.NewClient(paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                clientPingInterval, // Интервал ping'ов для обнаружения мёртвых соединений
			Timeout:             clientPingTimeout,  // Таймаут ожидания pong
			PermitWithoutStream: true,               // Держать соединение "тёплым" без активных RPC
		}))
	if err != nil {
		slog.Error("не удалось подключиться к PaymentService", "error", err)
		return
	}
	defer paymentConn.Close()

	// Создаём хранилище и обработчик
	store := ordstore.NewOrderStore()

	inventoryClient := inventory.NewInventoryClient(inventoryv1.NewInventoryServiceClient(inventoryConn))

	paymentClient := payment.NewPaymentClient(paymentv1.NewPaymentServiceClient(paymentConn))

	orderService := order.NewService(store, inventoryClient, paymentClient)

	orderApi := orderapi.NewApi(orderService)

	r, err := orderApi.SetupRouter()
	if err != nil {
		slog.Error("Не удалось инициализировать роутер", "error", err)
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
		serveErr := server.ListenAndServe()
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			slog.Error("ошибка запуска сервера", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	slog.Info("⚠️ Получен сигнал закрытия сервера. Выполняем graceful shutdown")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if shutdownErr := server.Shutdown(ctx); shutdownErr != nil {
		slog.Error("❌ ошибка при остановке сервера", "error", shutdownErr)
	}

	slog.Info("✅ Сервер остановлен")
}
