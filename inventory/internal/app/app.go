package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	"buf.build/go/protovalidate"
	protovalidateMiddleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/config"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/closer"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/grpc/health"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/logger"
	"github.com/PabloGolobaro/cosmic_factory/shared/pkg/interceptors"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

type App struct {
	diContainer *diContainer
	conf        config.Config
	listener    net.Listener
	grpcServer  *grpc.Server
}

// New создаёт и инициализирует приложение. Возвращает ошибку, если инициализация не удалась.
func New(ctx context.Context, conf config.Config) (*App, error) {
	a := &App{conf: conf}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

// Run управляет жизненным циклом приложения: запускает gRPC-сервер, обрабатывает сигналы ОС
// и выполняет graceful shutdown.
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a.startGracefulShutdown(ctx, cancel)

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initListener,
		a.initGRPCServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = newDIContainer(a.conf)
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.Init(a.conf.Logger.Level)
	return nil
}

func (a *App) initListener(_ context.Context) error {
	lis, err := new(net.ListenConfig).Listen(context.Background(), "tcp", a.conf.GRPC.Address())
	if err != nil {
		return fmt.Errorf("создание TCP-листенера: %w", err)
	}

	a.listener = lis
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	handler, err := a.diContainer.InventoryHandler(ctx)
	if err != nil {
		return err
	}

	validator, err := protovalidate.New()
	if err != nil {
		return fmt.Errorf("создание protovalidate валидатора: %w", err)
	}

	a.grpcServer = grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     a.conf.GRPC.MaxConnectionIdle,
			MaxConnectionAge:      a.conf.GRPC.MaxConnectionAge,
			MaxConnectionAgeGrace: a.conf.GRPC.MaxConnectionAgeGrace,
			Time:                  a.conf.GRPC.KeepaliveTime,
			Timeout:               a.conf.GRPC.KeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             a.conf.GRPC.MinPingInterval,
			PermitWithoutStream: true,
		}),
		grpc.ChainUnaryInterceptor(
			interceptors.RecoveryInterceptor(),
			interceptors.LoggerInterceptor(),
			protovalidateMiddleware.UnaryServerInterceptor(validator),
		),
	)

	inventoryv1.RegisterInventoryServiceServer(a.grpcServer, handler)
	health.RegisterService(a.grpcServer)
	reflection.Register(a.grpcServer)

	closer.Add("gRPC server", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	return nil
}

func (a *App) runGRPCServer() error {
	slog.Info("запуск InventoryService", "addr", a.conf.GRPC.Address())
	return a.grpcServer.Serve(a.listener)
}

// startGracefulShutdown запускает горутину, которая ожидает сигнал завершения (SIGINT/SIGTERM)
// и корректно останавливает все компоненты приложения через closer.
//
// Порядок работы:
//  1. Ждём сигнал ОС — signal.NotifyContext перехватывает его и отменяет ctx.
//  2. cancel() снимает перехват сигналов (вызывает signal.Stop внутри).
//     Контекст к этому моменту уже отменён. Это нужно, чтобы повторный Ctrl+C
//     не перехватывался, а завершил процесс принудительно (поведение ОС по умолчанию).
//  3. Создаём отдельный shutdownCtx от context.Background() с таймаутом — он не связан
//     с уже отменённым ctx, чтобы closer'ы имели гарантированное время на завершение.
func (a *App) startGracefulShutdown(ctx context.Context, cancel context.CancelFunc) {
	go func() { //nolint:gosec // G118: ctx уже отменён, context.Background нужен для graceful shutdown.
		<-ctx.Done()

		cancel()

		slog.Info("получен сигнал завершения, начинаем graceful shutdown")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), a.conf.GRPC.ShutdownTimeout)
		defer shutdownCancel()

		if closeErr := closer.CloseAll(shutdownCtx); closeErr != nil {
			slog.Error("ошибка при завершении работы", "error", closeErr)
		}
	}()
}
