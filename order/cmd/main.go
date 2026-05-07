package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/app"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("ошибка запуска сервера", "error", err)
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load("./../order.env"); err != nil {
		return fmt.Errorf("загрузка .env: %w", err)
	}

	configPath := config.ResolveConfigPath()

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("загрузка конфига: %w", err)
	}

	slog.Info("конфигурация загружена",
		"config_path", configPath,
		"http_address", cfg.HTTP.Address(),
		"grpc_address", cfg.GRPC.Address(),
		"pg_host", cfg.PG.Host,
		"inventory_address", cfg.Inventory.Address(),
		"payment_address", cfg.Payment.Address(),
	)

	a, err := app.New(context.Background(), *cfg)
	if err != nil {
		return fmt.Errorf("инициализация приложения: %w", err)
	}

	return a.Run()
}
