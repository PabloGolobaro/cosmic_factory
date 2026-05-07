package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Logger loggerConfig `yaml:"logger"`
	GRPC   grpcConfig   `yaml:"grpc"`
	PG     pgConfig     `yaml:"pg"`
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
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			return nil, fmt.Errorf("не удалось загрузить конфиг из %q: %w", path, err)
		}

		return &cfg, nil
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("не удалось загрузить конфиг из env: %w", err)
	}

	return &cfg, nil
}
