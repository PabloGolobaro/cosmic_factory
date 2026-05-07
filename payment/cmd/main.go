package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/PabloGolobaro/cosmic_factory/payment/internal/app"
	"github.com/PabloGolobaro/cosmic_factory/payment/internal/config"
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
		return fmt.Errorf("загрузка конфига: %w", err)
	}

	slog.Info("конфигурация загружена",
		"config_path", configPath,
		"grpc_address", cfg.GRPC.Address(),
	)

	a, err := app.New(context.Background(), *cfg)
	if err != nil {
		return fmt.Errorf("инициализация приложения: %w", err)
	}

	return a.Run()
}
