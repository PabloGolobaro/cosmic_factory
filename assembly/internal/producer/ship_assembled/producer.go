package shipassembled

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/PabloGolobaro/cosmic_factory/assembly/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka"
	eventsv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/events/v1"
)

type service struct {
	producer KafkaProducer
}

func NewService(producer KafkaProducer) *service {
	return &service{producer: producer}
}

func (s *service) PublishShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	pb := &eventsv1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
		AssembledAt:  timestamppb.New(event.AssembledAt),
	}

	data, err := proto.Marshal(pb)
	if err != nil {
		return fmt.Errorf("сериализация ShipAssembled: %w", err)
	}

	return s.producer.Send(ctx, &kafka.Message{
		Key:       []byte(event.OrderUUID),
		Value:     data,
		Timestamp: time.Now(),
	})
}
