package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Logger    loggerConfig          `yaml:"logger"`
	OTel      otelConfig            `yaml:"otel"`
	HTTP      httpConfig            `yaml:"http"`
	PG        pgConfig              `yaml:"pg"`
	Inventory inventoryClientConfig `yaml:"inventory"`
	Payment   paymentClientConfig   `yaml:"payment"`
	IAM       iamClientConfig       `yaml:"iam"`
	Kafka     kafkaConfig           `yaml:"kafka"`
	RateLimit rateLimitConfig       `yaml:"rate_limit"`
}

// LogValue реализует slog.LogValuer — все поля без пароля.
func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Group("logger",
			slog.String("level", c.Logger.Level),
		),
		slog.Group("http",
			slog.String("host", c.HTTP.Host),
			slog.String("port", c.HTTP.Port),
			slog.Duration("read_header_timeout", c.HTTP.ReadHeaderTimeout),
			slog.Duration("read_timeout", c.HTTP.ReadTimeout),
			slog.Duration("write_timeout", c.HTTP.WriteTimeout),
			slog.Duration("idle_timeout", c.HTTP.IdleTimeout),
			slog.Duration("shutdown_timeout", c.HTTP.ShutdownTimeout),
		),
		slog.Group("pg",
			slog.String("host", c.PG.Host),
			slog.String("port", c.PG.Port),
			slog.String("database", c.PG.Database),
			slog.String("user", c.PG.User),
			slog.String("password", "***"),
			slog.String("sslmode", c.PG.SSLMode),
		),
		slog.Group("otel",
			slog.String("endpoint", c.OTel.Endpoint),
			slog.String("service_name", c.OTel.ServiceName),
		),
		slog.Group("inventory",
			slog.String("host", c.Inventory.Host),
			slog.String("port", c.Inventory.Port),
			slog.Duration("ping_interval", c.Inventory.PingInterval),
			slog.Duration("ping_timeout", c.Inventory.PingTimeout),
		),
		slog.Group("payment",
			slog.String("host", c.Payment.Host),
			slog.String("port", c.Payment.Port),
			slog.Duration("ping_interval", c.Payment.PingInterval),
			slog.Duration("ping_timeout", c.Payment.PingTimeout),
		),
		slog.Group("iam",
			slog.String("host", c.IAM.Host),
			slog.String("port", c.IAM.Port),
			slog.Duration("ping_interval", c.IAM.PingInterval),
			slog.Duration("ping_timeout", c.IAM.PingTimeout),
		),
		slog.Group("kafka",
			slog.Any("brokers", c.Kafka.Brokers),
			slog.String("group", c.Kafka.ConsumerGroup),
			slog.String("consume", c.Kafka.ConsumeTopic),
			slog.String("produce", c.Kafka.ProduceTopic),
		),
		slog.Group("rate_limit",
			slog.String("redis", c.RateLimit.RedisAddress),
			slog.Int("rate", c.RateLimit.Rate),
			slog.Int("burst", c.RateLimit.Burst),
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
