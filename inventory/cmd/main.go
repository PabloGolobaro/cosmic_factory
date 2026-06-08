package main

import (
	"context"
	"errors"
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
	if err := godotenv.Load("./../inventory.env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("загрузка .env: %w", err)
	}

	configPath, configSource := config.ResolveConfigPath()

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("загрузка конфига: %w", err)
	}

	slog.Info("конфигурация загружена",
		"source", configSource,
		"path", configPath,
		"config", cfg,
	)

	a, err := app.New(context.Background(), *cfg)
	if err != nil {
		return fmt.Errorf("инициализация приложения: %w", err)
	}

	return a.Run()
}
