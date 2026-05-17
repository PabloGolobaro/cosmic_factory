package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"

	"github.com/PabloGolobaro/cosmic_factory/assembly/internal/config"
	orderpaid "github.com/PabloGolobaro/cosmic_factory/assembly/internal/consumer/order_paid"
	shipassembled "github.com/PabloGolobaro/cosmic_factory/assembly/internal/producer/ship_assembled"
	assemblyservice "github.com/PabloGolobaro/cosmic_factory/assembly/internal/service/assembly"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/closer"
	kafkaconsumer "github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka/consumer"
	kafkaproducer "github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka/producer"
	kafkamw "github.com/PabloGolobaro/cosmic_factory/platform/pkg/middleware/kafka"
)

type consumerRunner interface {
	RunConsumer(ctx context.Context) error
}

type diContainer struct {
	conf config.Config

	consumerGroup sarama.ConsumerGroup
	syncProducer  sarama.SyncProducer

	shipProducer    assemblyservice.Producer
	assemblySvc     orderpaid.AssemblyService
	orderPaidRunner consumerRunner
}

func newDIContainer(conf config.Config) *diContainer {
	return &diContainer{conf: conf}
}

func (d *diContainer) KafkaConsumerGroup() (sarama.ConsumerGroup, error) {
	if d.consumerGroup == nil {
		cfg := sarama.NewConfig()
		cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
		cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
			sarama.NewBalanceStrategyRoundRobin(),
		}

		group, err := sarama.NewConsumerGroup(d.conf.Kafka.Brokers, d.conf.Kafka.ConsumerGroup, cfg)
		if err != nil {
			return nil, fmt.Errorf("создание Kafka consumer group: %w", err)
		}

		slog.Info("Kafka consumer group создана", "group", d.conf.Kafka.ConsumerGroup)

		closer.Add("kafka consumer group", func(_ context.Context) error {
			return group.Close()
		})

		d.consumerGroup = group
	}

	return d.consumerGroup, nil
}

func (d *diContainer) KafkaSyncProducer() (sarama.SyncProducer, error) {
	if d.syncProducer == nil {
		cfg := sarama.NewConfig()
		cfg.Producer.Return.Successes = true
		cfg.Producer.RequiredAcks = sarama.WaitForAll

		producer, err := sarama.NewSyncProducer(d.conf.Kafka.Brokers, cfg)
		if err != nil {
			return nil, fmt.Errorf("создание Kafka sync producer: %w", err)
		}

		slog.Info("Kafka sync producer создан", "topic", d.conf.Kafka.ProduceTopic)

		closer.Add("kafka sync producer", func(_ context.Context) error {
			return producer.Close()
		})

		d.syncProducer = producer
	}

	return d.syncProducer, nil
}

func (d *diContainer) ShipAssembledProducer() (assemblyservice.Producer, error) {
	if d.shipProducer == nil {
		syncProducer, err := d.KafkaSyncProducer()
		if err != nil {
			return nil, fmt.Errorf("ship assembled producer: %w", err)
		}

		p := kafkaproducer.NewProducer(syncProducer, d.conf.Kafka.ProduceTopic)
		d.shipProducer = shipassembled.NewService(p)
	}

	return d.shipProducer, nil
}

func (d *diContainer) AssemblyService() (orderpaid.AssemblyService, error) {
	if d.assemblySvc == nil {
		producer, err := d.ShipAssembledProducer()
		if err != nil {
			return nil, fmt.Errorf("assembly service: %w", err)
		}

		d.assemblySvc = assemblyservice.NewService(producer)
	}

	return d.assemblySvc, nil
}

func (d *diContainer) OrderPaidConsumerService() (consumerRunner, error) {
	if d.orderPaidRunner == nil {
		group, err := d.KafkaConsumerGroup()
		if err != nil {
			return nil, fmt.Errorf("order paid consumer: %w", err)
		}

		assemblySvc, err := d.AssemblyService()
		if err != nil {
			return nil, fmt.Errorf("order paid consumer: %w", err)
		}

		consumer := kafkaconsumer.NewConsumer(
			group,
			[]string{d.conf.Kafka.ConsumeTopic},
			kafkaconsumer.WithMiddlewares(kafkamw.ConsumerLogging()),
		)

		d.orderPaidRunner = orderpaid.NewService(consumer, assemblySvc)
	}

	return d.orderPaidRunner, nil
}
