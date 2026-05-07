package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/app"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("ошибка запуска сервера", "error", err)
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load("./../inventory.env"); err != nil {
		return fmt.Errorf("загрузка .env: %w", err)
	}

	configPath := config.ResolveConfigPath()

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("загрузка конфига: %w", err)
	}

	slog.Info("конфигурация загружена",
		"config_path", configPath,
		"grpc_address", cfg.GRPC.Address(),
		"pg_host", cfg.PG.Host,
	)

	a, err := app.New(context.Background(), *cfg)
	if err != nil {
		return fmt.Errorf("инициализация приложения: %w", err)
	}

	return a.Run()
}
