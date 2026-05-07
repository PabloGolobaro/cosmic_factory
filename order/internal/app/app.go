package app

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/go-faster/errors"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/config"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/closer"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
	conf        config.Config
	server      *http.Server
}

// New создаёт и инициализирует приложение. Возвращает ошибку, если инициализация не удалась.
func New(ctx context.Context, conf config.Config) (*App, error) {
	a := &App{conf: conf}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

// Run управляет жизненным циклом приложения: запускает HTTP-сервер, обрабатывает сигналы ОС
// и выполняет graceful shutdown.
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a.startGracefulShutdown(ctx, cancel)

	return a.runHTTPServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initHTTPServer,
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

func (a *App) initHTTPServer(ctx context.Context) error {
	r, err := a.diContainer.Router(ctx)
	if err != nil {
		return err
	}

	a.server = &http.Server{
		Addr:              a.conf.HTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: a.conf.HTTP.ReadHeaderTimeout,
		ReadTimeout:       a.conf.HTTP.ReadTimeout,
		WriteTimeout:      a.conf.HTTP.WriteTimeout,
		IdleTimeout:       a.conf.HTTP.IdleTimeout,
	}

	closer.Add("HTTP server", func(shutdownCtx context.Context) error {
		return a.server.Shutdown(shutdownCtx)
	})

	return nil
}

func (a *App) runHTTPServer() error {
	slog.Info("запуск OrderService", "addr", a.conf.HTTP.Address())
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
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

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), a.conf.HTTP.ShutdownTimeout)
		defer shutdownCancel()

		if closeErr := closer.CloseAll(shutdownCtx); closeErr != nil {
			slog.Error("ошибка при завершении работы", "error", closeErr)
		}
	}()
}
