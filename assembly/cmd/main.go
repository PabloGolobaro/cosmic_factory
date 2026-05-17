package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/PabloGolobaro/cosmic_factory/assembly/internal/app"
	"github.com/PabloGolobaro/cosmic_factory/assembly/internal/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("ошибка запуска сервиса", "error", err)
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
		"consume_topic", cfg.Kafka.ConsumeTopic,
		"produce_topic", cfg.Kafka.ProduceTopic,
	)

	a, err := app.New(context.Background(), *cfg)
	if err != nil {
		return fmt.Errorf("инициализация приложения: %w", err)
	}

	return a.Run()
}
