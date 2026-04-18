package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"buf.build/go/protovalidate"
	protovalidateMiddleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	svc "github.com/PabloGolobaro/cosmic_factory/inventory/pkg/service"
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
	listener := &net.ListenConfig{}

	lis, err := listener.Listen(context.Background(), "tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		os.Exit(1)
	}

	// Создаем protovalidate валидатор для проверки входящих запросов
	validator, err := protovalidate.New()
	if err != nil {
		slog.Error("ошибка создания валидатора", "error", err)
		return
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

	// Интерцепторы: recovery (перехват паник) + логирование запросов

	inventoryv1.RegisterInventoryServiceServer(grpcServer, svc.NewInventoryServer())

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	go func() {
		slog.Info("запуск InventoryService", "адрес", grpcAddress)
		err = grpcServer.Serve(lis)
		if err != nil {
			slog.Error("ошибка запуска сервера", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	slog.Info("⚠️ Получен сигнал закрытия сервера. Выполняем graceful shutdown")

	grpcServer.GracefulStop()

	slog.Info("✅ Сервер остановлен")
}
