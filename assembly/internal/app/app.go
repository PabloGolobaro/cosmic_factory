package app

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/PabloGolobaro/cosmic_factory/assembly/internal/config"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/closer"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
	conf        config.Config
}

func New(ctx context.Context, conf config.Config) (*App, error) {
	a := &App{conf: conf}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a.startGracefulShutdown(ctx, cancel)

	return a.runConsumer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
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

func (a *App) runConsumer(ctx context.Context) error {
	slog.InfoContext(ctx, "запуск AssemblyService")

	svc, err := a.diContainer.OrderPaidConsumerService()
	if err != nil {
		return err
	}

	return svc.RunConsumer(ctx)
}

func (a *App) startGracefulShutdown(ctx context.Context, cancel context.CancelFunc) {
	go func() { //nolint:gosec // G118: ctx уже отменён, context.Background нужен для graceful shutdown.
		<-ctx.Done()

		cancel()

		slog.Info("получен сигнал завершения, начинаем graceful shutdown")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), a.conf.App.ShutdownTimeout)
		defer shutdownCancel()

		if closeErr := closer.CloseAll(shutdownCtx); closeErr != nil {
			slog.Error("ошибка при завершении работы", "error", closeErr)
		}
	}()
}
