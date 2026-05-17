package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Logger    loggerConfig          `yaml:"logger"`
	GRPC      grpcConfig            `yaml:"grpc"`
	HTTP      httpConfig            `yaml:"http"`
	PG        pgConfig              `yaml:"pg"`
	Inventory inventoryClientConfig `yaml:"inventory"`
	Payment   paymentClientConfig   `yaml:"payment"`
	Kafka     kafkaConfig           `yaml:"kafka"`
}

const defaultConfigPath = "config.local.yaml"

// ResolveConfigPath определяет путь к конфиг-файлу по цепочке приоритетов:
// флаг -config > env CONFIG_PATH > "config.local.yaml".
func ResolveConfigPath() string {
	var cfgFlag string
	flag.StringVar(&cfgFlag, "config", "", "путь к YAML-конфигу (например, config.staging.yaml)")
	flag.Parse()

	if cfgFlag != "" {
		return cfgFlag
	}

	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		return envPath
	}

	return defaultConfigPath
}

func Load(path string) (*Config, error) {
	var cfg Config

	if path != "" {
		// ReadConfig читает YAML-файл, а затем перетирает значения из env-переменных.
		// Приоритет: env > yaml > env-default.
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			return nil, fmt.Errorf("не удалось загрузить конфиг из %q: %w", path, err)
		}

		return &cfg, nil
	}

	// Если путь не указан — читаем только из env-переменных.
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("не удалось загрузить конфиг из env: %w", err)
	}

	return &cfg, nil
}
