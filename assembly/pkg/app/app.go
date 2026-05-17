package app

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/IBM/sarama"

	orderpaid "github.com/PabloGolobaro/cosmic_factory/assembly/internal/consumer/order_paid"
	shipassembled "github.com/PabloGolobaro/cosmic_factory/assembly/internal/producer/ship_assembled"
	assemblyservice "github.com/PabloGolobaro/cosmic_factory/assembly/internal/service/assembly"
	kafkaconsumer "github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka/consumer"
	kafkaproducer "github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka/producer"
)

// Config — лёгкая конфигурация для e2e-тестов: принимает уже созданные Sarama-объекты.
type Config struct {
	OrderPaidTopic     string
	ShipAssembledTopic string
	MinBuildTimeSec    int64
	MaxBuildTimeSec    int64
}

// App — собранный assembly-сервис, готовый к запуску.
type App struct {
	runner interface {
		RunConsumer(ctx context.Context) error
	}
}

// New собирает assembly-сервис из готовых Sarama producer/consumer group.
func New(syncProducer sarama.SyncProducer, cg sarama.ConsumerGroup, cfg Config) *App {
	p := kafkaproducer.NewProducer(syncProducer, cfg.ShipAssembledTopic)
	shipProducer := shipassembled.NewService(p)

	svc := assemblyservice.NewService(shipProducer)
	svc.SetBuildDelay(makeBuildDelay(cfg.MinBuildTimeSec, cfg.MaxBuildTimeSec))

	consumer := kafkaconsumer.NewConsumer(cg, []string{cfg.OrderPaidTopic})
	runner := orderpaid.NewService(consumer, svc)

	return &App{runner: runner}
}

func (a *App) RunConsumer(ctx context.Context) error {
	return a.runner.RunConsumer(ctx)
}

func makeBuildDelay(minSec, maxSec int64) func() time.Duration {
	if maxSec <= minSec {
		return func() time.Duration { return time.Duration(minSec) * time.Second }
	}

	spread := maxSec - minSec + 1
	return func() time.Duration {
		return time.Duration(minSec+rand.Int64N(spread)) * time.Second //nolint:mnd
	}
}
