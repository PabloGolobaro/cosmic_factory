package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"buf.build/go/protovalidate"
	protovalidateMiddleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	v1 "github.com/PabloGolobaro/cosmic_factory/payment/internal/api/payment/v1"
	"github.com/PabloGolobaro/cosmic_factory/payment/internal/config"
	"github.com/PabloGolobaro/cosmic_factory/payment/internal/service/payment"
	"github.com/PabloGolobaro/cosmic_factory/shared/pkg/interceptors"
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

func main() {
	if err := run(); err != nil {
		slog.Error("ошибка запуска сервера", "error", err)
		os.Exit(1)
	}
}

func run() error {
	configPath := config.ResolveConfigPath()

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	slog.Info("конфигурация загружена",
		"config_path", configPath,
		"grpc_address", cfg.GRPC.Address(),
	)

	lis, err := new(net.ListenConfig).Listen(context.Background(), "tcp", cfg.GRPC.Address())
	if err != nil {
		return fmt.Errorf("создание listener: %w", err)
	}

	validator, err := protovalidate.New()
	if err != nil {
		return fmt.Errorf("создание валидатора: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     cfg.GRPC.MaxConnectionIdle,
			MaxConnectionAge:      cfg.GRPC.MaxConnectionAge,
			MaxConnectionAgeGrace: cfg.GRPC.MaxConnectionAgeGrace,
			Time:                  cfg.GRPC.KeepaliveTime,
			Timeout:               cfg.GRPC.KeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             cfg.GRPC.MinPingInterval,
			PermitWithoutStream: true,
		}),
		grpc.ChainUnaryInterceptor(
			interceptors.RecoveryInterceptor(),
			interceptors.LoggerInterceptor(),
			protovalidateMiddleware.UnaryServerInterceptor(validator),
		),
	)

	svc := payment.NewPaymentService()
	api := v1.NewApi(svc)
	paymentv1.RegisterPaymentServiceServer(grpcServer, api)
	reflection.Register(grpcServer)

	go func() {
		slog.Info("запуск PaymentService", "адрес", cfg.GRPC.Address())
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
