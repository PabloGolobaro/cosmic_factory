package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Logger  loggerConfig  `yaml:"logger"`
	OTel    otelConfig    `yaml:"otel"`
	GRPC    grpcConfig    `yaml:"grpc"`
	PG      pgConfig      `yaml:"pg"`
	Redis   redisConfig   `yaml:"redis"`
	Session sessionConfig `yaml:"session"`
}

// LogValue реализует slog.LogValuer — все поля без пароля.
func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Group("logger",
			slog.String("level", c.Logger.Level),
		),
		slog.Group("grpc",
			slog.String("host", c.GRPC.Host),
			slog.String("port", c.GRPC.Port),
			slog.Duration("max_connection_idle", c.GRPC.MaxConnectionIdle),
			slog.Duration("max_connection_age", c.GRPC.MaxConnectionAge),
			slog.Duration("max_connection_age_grace", c.GRPC.MaxConnectionAgeGrace),
			slog.Duration("keepalive_time", c.GRPC.KeepaliveTime),
			slog.Duration("keepalive_timeout", c.GRPC.KeepaliveTimeout),
			slog.Duration("min_ping_interval", c.GRPC.MinPingInterval),
			slog.Duration("shutdown_timeout", c.GRPC.ShutdownTimeout),
		),
		slog.Group("pg",
			slog.String("host", c.PG.Host),
			slog.String("port", c.PG.Port),
			slog.String("database", c.PG.Database),
			slog.String("user", c.PG.User),
			slog.String("password", "***"),
			slog.String("sslmode", c.PG.SSLMode),
		),
		slog.Group("redis",
			slog.String("host", c.Redis.Host),
			slog.String("port", c.Redis.Port),
		),
		slog.Group("otel",
			slog.String("endpoint", c.OTel.Endpoint),
			slog.String("service_name", c.OTel.ServiceName),
		),
		slog.Group("session",
			slog.Duration("ttl", c.Session.TTL),
		),
	)
}

const defaultConfigPath = "config.local.yaml"

// ResolveConfigPath определяет путь к конфиг-файлу по цепочке приоритетов:
// флаг -config > env CONFIG_PATH > "config.local.yaml".
// Возвращает путь и источник: "flag" | "env" | "default".
func ResolveConfigPath() (path, source string) {
	var cfgFlag string
	flag.StringVar(&cfgFlag, "config", "", "путь к YAML-конфигу (например, config.staging.yaml)")
	flag.Parse()

	if cfgFlag != "" {
		return cfgFlag, "flag"
	}

	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		return envPath, "env"
	}

	return defaultConfigPath, "default"
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
